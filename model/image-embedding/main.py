from flask import Flask, request, jsonify
import torch
from transformers import CLIPProcessor, CLIPModel
from PIL import Image
import torchvision.transforms as transforms
from surprise import Dataset, Reader, SVD, Trainset
import pandas as pd
from flask_cors import CORS
import numpy as np
import logging
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.metrics.pairwise import linear_kernel
import pandas as pd

logging.basicConfig(level=logging.INFO,
                    format='%(asctime)s %(levelname)s %(message)s',
                    handlers=[logging.StreamHandler()])
logger = logging.getLogger(__name__)

app = Flask(__name__)
CORS(app)

# CLIP model setup
device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
model_name = "openai/clip-vit-base-patch32"
clip = CLIPModel.from_pretrained(model_name).to(device)
processor = CLIPProcessor.from_pretrained(model_name)

# SVD model setup
ratings_file_path = '../../data/movielens/ratings.csv'
reader = Reader(line_format='user item rating timestamp', sep=',', rating_scale=(1, 5), skip_lines=1)
data = Dataset.load_from_file(ratings_file_path, reader=reader)
trainset = data.build_full_trainset()
movies_file_path = '../../data/movielens/movies.csv'
links_file_path = '../../data/movielens/links.csv'
movies = pd.read_csv(movies_file_path)
links = pd.read_csv(links_file_path).dropna(subset=['tmdbId'])  # Drop rows with NA in 'tmdbId'
algo = SVD(n_factors=100, reg_all=0.02, lr_all=0.005) 
algo.fit(trainset)

# Dictionary to track clicked items for each user
clicked_items_tracker = {}

# Mapping from tmdbId to movieId, ensuring tmdbId is treated as an integer
tmdb_to_movieId = pd.Series(links['movieId'].values, index=links['tmdbId'].astype(int)).to_dict()
movieId_to_tmdb = pd.Series(links['tmdbId'].values, index=links['movieId'].astype(int)).to_dict()


ratings_df = pd.read_csv(ratings_file_path)
movies_df = pd.read_csv(movies_file_path)
links_df = pd.read_csv(links_file_path)
movies_with_ratings = pd.merge(ratings_df, movies_df, on='movieId')
tfidf = TfidfVectorizer(stop_words='english')
tfidf_matrix = tfidf.fit_transform(movies_df['genres'])
cosine_sim = linear_kernel(tfidf_matrix, tfidf_matrix)

def get_recommendations(title, cosine_sim=cosine_sim):
    idx = movies_df.loc[movies_df['title'] == title].index[0]
    sim_scores = list(enumerate(cosine_sim[idx]))
    sim_scores = sorted(sim_scores, key=lambda x: x[1], reverse=True)
    sim_scores = sim_scores[1:2]  # Get the top 10 most similar movies
    movie_indices = [i[0] for i in sim_scores]
    return movies_df['movieId'].iloc[movie_indices].tolist()

@app.route('/encode', methods=['POST'])
def encode_text():
    data = request.get_json()
    text = data['text']
    print(text)
    with torch.no_grad():
        inputs = processor(text, return_tensors="pt", padding=True, truncation=True).to(device)
        text_features = clip.get_text_features(**inputs)
    resp = jsonify({"text_features": text_features.cpu().numpy().tolist()[0]})
    print(text_features.cpu().numpy().tolist()[0])
    return resp

@app.route('/predict/<int:user_id>', methods=['GET'])
def predict(user_id):
    try:
        # Predict ratings for all movies
        all_movies = movies['movieId'].unique()
        predictions = [(movie_id, algo.predict(user_id, movie_id).est) for movie_id in all_movies]
        logger.info(predictions)

        # Convert ratings to binary (1 for liked, 0 for not liked)
        liked_movies = [movie_id for movie_id, rating in predictions if rating >= 4]

        # Get movie titles for the liked movies
        liked_movie_titles = movies[movies['movieId'].isin(liked_movies)]['title']

        return jsonify({"user_id": user_id, "liked_movies": liked_movie_titles.tolist()})
    except ValueError as e:
        return jsonify({"error": str(e)}), 400

def update_trainset(user_id, clicked_items):
    global trainset
    # Convert current trainset to dataframe
    ratings_df = pd.DataFrame(trainset.all_ratings(), columns=['user', 'item', 'rating'])
    # Convert inner ids to raw ids
    ratings_df['user'] = ratings_df['user'].apply(trainset.to_raw_uid)
    ratings_df['item'] = ratings_df['item'].apply(trainset.to_raw_iid)
    
    # Ensure clicked_items are treated as integers and map tmdbId to movieId
    # clicked_movie_ids = [tmdb_to_movieId[int(tmdb_id)] for tmdb_id in clicked_items if int(tmdb_id) in tmdb_to_movieId]

    # Add new user ratings
    new_ratings = pd.DataFrame([{'user': user_id, 'item': item_id, 'rating': 5.0} for item_id in clicked_items])  # Assuming binary rating
    logger.info(new_ratings)
    ratings_df = pd.concat([ratings_df, new_ratings], ignore_index=True)
    
    # Load new ratings into the trainset
    reader = Reader(rating_scale=(1, 5))
    data = Dataset.load_from_df(ratings_df[['user', 'item', 'rating']], reader)
    trainset = data.build_full_trainset()

@app.route('/train_model/<int:user_id>', methods=['POST'])
def train_model(user_id):
    try:
        # Get clicked items from request
        clicked_items = request.json.get('clicked_items')
        if not clicked_items:
            return jsonify({"error": "No clicked items provided"}), 400

        # Track clicked items for the user
        if user_id not in clicked_items_tracker:
            clicked_items_tracker[user_id] = []

        # Ensure clicked_items are treated as integers and map tmdbId to movieId
        clicked_movie_ids = [tmdb_to_movieId[int(tmdb_id)] for tmdb_id in clicked_items if int(tmdb_id) in tmdb_to_movieId]
        clicked_movies = movies[movies['movieId'].isin(clicked_movie_ids)]['title']
        clicked_items_tracker[user_id].extend(clicked_movies)

        logger.info(clicked_items_tracker)
        logger.info({"user_id": user_id, "all clicked_movies": clicked_items_tracker[user_id]})

        # Update trainset with new user ratings
        update_trainset(user_id, clicked_movie_ids)

        # Retrain the model
        algo.fit(trainset)
        logger.info(trainset)

        return jsonify({"message": "Model retrained successfully"})
    except ValueError as e:
        return jsonify({"error": str(e)}), 400


@app.route('/click/<int:user_id>', methods=['POST'])
def click(user_id):
    try:
        # Get clicked items from request
        clicked_items = request.json.get('clicked_items')
        if not clicked_items:
            return jsonify({"error": "No clicked items provided"}), 400

        # Track clicked items for the user
        if user_id not in clicked_items_tracker:
            clicked_items_tracker[user_id] = []

        # Ensure clicked_items are treated as integers and map tmdbId to movieId
        clicked_movie_ids = [tmdb_to_movieId[int(tmdb_id)] for tmdb_id in clicked_items if int(tmdb_id) in tmdb_to_movieId]
        clicked_movies = movies[movies['movieId'].isin(clicked_movie_ids)]['title']
        clicked_items_tracker[user_id].extend(clicked_movies)

        logger.info(clicked_items_tracker)
        logger.info({"user_id": user_id, "all clicked_movies": clicked_items_tracker[user_id]})

        return jsonify({"message": "Model retrained successfully"})
    except ValueError as e:
        return jsonify({"error": str(e)}), 400

# @app.route('/similar_items/<int:user_id>', methods=['GET'])
# def get_similar_items(user_id):
#     try: 
#         similar_items = []
#         movie_titles = clicked_items_tracker[user_id]
#         for movie_title in movie_titles:
#             similar_item = get_recommendations(movie_title)
#             print(similar_item)
#             tmdbId = links_df.loc[links_df['movieId'] == similar_item[0], 'tmdbId'].values[0]
#             print(tmdbId)
#             similar_items.append(int(tmdbId))
#         return jsonify({"ids": similar_items})
#     except ValueError as e:
#         return jsonify({"error": str(e)}), 

@app.route('/similar_items/<int:user_id>', methods=['GET'])
def get_similar_items(user_id):
    try:
        similar_items_set = set()  # Use a set to store unique tmdbIds
        movie_titles = clicked_items_tracker[user_id]
        for movie_title in movie_titles:
            similar_movies = get_recommendations(movie_title)
            print(similar_movies)
            for similar_movie in similar_movies:
                print("getting tmdb_id")
                tmdb_id = links_df.loc[links_df['movieId'] == similar_movie, 'tmdbId'].values[0]
                print("getting tmdb_id")
                similar_items_set.add(int(tmdb_id))  # Add the tmdbId to the set
        print(similar_items_set)
        similar_items = list(similar_items_set)  # Convert the set to a list
        print(similar_items)
        return jsonify({"ids": similar_items})
    except Exception as e:
        return jsonify({"error": str(e)}), 400

if __name__ == '__main__':
    app.run(debug=True)