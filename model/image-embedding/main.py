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
import torch
import torch.nn as nn
import torch.optim as optim
from torch.utils.data import DataLoader, TensorDataset
import pandas as pd
import json
from llamaapi import LlamaAPI
from retry import retry

from sklearn.model_selection import train_test_split

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

############ MF model
class MatrixFactorization(nn.Module):
    def __init__(self, num_users, num_items, embedding_dim):
        super(MatrixFactorization, self).__init__()
        self.user_embedding = nn.Embedding(num_users, embedding_dim)
        self.item_embedding = nn.Embedding(num_items, embedding_dim)

    def forward(self, user_ids, item_ids):
        user_embeds = self.user_embedding(user_ids)
        item_embeds = self.item_embedding(item_ids)
        return (user_embeds * item_embeds).sum(1)

ratings_data = pd.read_csv('../../data/movielens/ratings.csv')
movies_data = pd.read_csv('../../data/movielens/movies.csv')

user_item_matrix = ratings_data.pivot(index='userId', columns='movieId', values='rating').fillna(0)

user_item_tensor = torch.tensor(user_item_matrix.values, dtype=torch.float32)

# Split the user-item interaction matrix into train and test sets
train_matrix, test_matrix = train_test_split(user_item_tensor, test_size=0.2, random_state=42)

# Data Loaders for train and test sets
train_loader = DataLoader(TensorDataset(*torch.where(train_matrix != 0)), batch_size=64, shuffle=True)
test_loader = DataLoader(TensorDataset(*torch.where(test_matrix != 0)), batch_size=64, shuffle=False)

# Model, Loss, Optimizer
num_users, num_items = user_item_tensor.size()
embedding_dim = 50
lr = 0.01
epochs = 50
model = MatrixFactorization(num_users, num_items, embedding_dim)
criterion = nn.MSELoss()
optimizer = optim.SGD(model.parameters(), lr=lr)

def train_model(user_item_matrix, epochs=50):
    user_item_tensor = torch.tensor(user_item_matrix.values, dtype=torch.float32)

    # Split the user-item interaction matrix into train and test sets
    train_matrix, test_matrix = train_test_split(user_item_tensor, test_size=0.2, random_state=42)

    # Data Loaders for train and test sets
    train_loader = DataLoader(TensorDataset(*torch.where(train_matrix != 0)), batch_size=64, shuffle=True)
    test_loader = DataLoader(TensorDataset(*torch.where(test_matrix != 0)), batch_size=64, shuffle=False)

    num_users, num_items = user_item_tensor.size()
    print("num_users training", num_users)
    embedding_dim = 50
    lr = 0.01

    for epoch in range(epochs):
        total_loss = 0.0
        total_samples = 0
        for batch in train_loader:
            optimizer.zero_grad()
            user_ids, item_ids = batch
            outputs = model(user_ids, item_ids)
            loss = criterion(outputs, train_matrix[user_ids, item_ids])
            loss.backward()
            optimizer.step()

            total_loss += loss.item() * len(user_ids)
            total_samples += len(user_ids)

        train_loss = total_loss / total_samples
        print(f"Epoch {epoch+1}, Training Loss: {train_loss:.4f}")

        # Evaluation
        with torch.no_grad():
            total_loss = 0.0
            total_samples = 0
            for batch in test_loader:
                user_ids, item_ids = batch
                outputs = model(user_ids, item_ids)
                loss = criterion(outputs, test_matrix[user_ids, item_ids])

                total_loss += loss.item() * len(user_ids)
                total_samples += len(user_ids)

            test_loss = total_loss / total_samples
            print(f"Epoch {epoch+1}, Test Loss: {test_loss:.4f}")

train_model(user_item_matrix)

# content based
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

@app.route('/get_recommendations/<int:user_id>', methods=['GET'])
def predict(user_id):
    try:
        global num_items, num_users
        
        user_ids = torch.full((num_items,), user_id, dtype=torch.long)
        item_ids = torch.arange(num_items, dtype=torch.long)
        print(num_items, user_ids)
        with torch.no_grad():
            predictions = model(user_ids, item_ids)
        top_movie_ids = torch.argsort(predictions, descending=True)[:10]
        liked_movie_titles = movies_data[movies_data['movieId'].isin(top_movie_ids.numpy())]['movieId']

        return jsonify({"user_id": user_id, "liked_movies": liked_movie_titles.tolist()})
    except ValueError as e:
        return jsonify({"error": str(e)}), 400


def update_user_ratings_with_clicks(user_id, clicked_movies, click_rating, user_item_matrix):
    print("before update", len(user_item_matrix.index))
    # Check if user exists in the user_item_matrix
    if user_id in user_item_matrix.index:
        # Update ratings for the clicked movies
        for movie_id in clicked_movies:
            if movie_id in user_item_matrix.columns:
                user_item_matrix.at[user_id, movie_id] = click_rating
    else:
        # Add new row for new user with ratings for clicked movies
        new_row = pd.Series(0, index=user_item_matrix.columns, name=user_id)
        for movie_id in clicked_movies:
            if movie_id in user_item_matrix.columns:
                new_row[movie_id] = click_rating
        user_item_matrix = pd.concat([user_item_matrix, new_row.to_frame().T])
    
    print("after update", len(user_item_matrix.index))

    return user_item_matrix

@app.route('/click/<int:user_id>', methods=['POST'])
def click(user_id):
    global user_item_matrix
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

        # user_item_matrix = update_user_ratings_with_clicks(user_id, clicked_movie_ids, 5, user_item_matrix)
        # train_model(user_item_matrix, epochs=10)

        return jsonify({"message": "Model retrained successfully"})
    except ValueError as e:
        return jsonify({"error": str(e)}), 400

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
                similar_items_set.add(int(tmdb_id))  # Add the tmdbId to the set
        print(similar_items_set)
        similar_items = list(similar_items_set)  # Convert the set to a list
        print(similar_items)
        return jsonify({"ids": similar_items})
    except Exception as e:
        return jsonify({"error": str(e)}), 400

llama = LlamaAPI("LL-5JPtGu2fBHzbE5JrkJsHwZmUrgn1shI1FAow28IeaTuqZvIke3oEKsLRMmTM9xPc")

@app.route('/get_nlp_keywords/<string:q>', methods=['GET'])
@retry(tries=3, delay=2, backoff=2)
def get_nlp_keywords(q):
    try:
        my_prompt = "Give me 3 keywords only, separated by commas"
        api_request_json = {
            "messages": [
                {"role": "user", "content": q+my_prompt},
            ]
        }
        response = llama.run(api_request_json)
        return jsonify({"keywords": response.json()['choices'][0]['message']['content']})
    except Exception as e:
        return jsonify({"error": str(e)}), 400


if __name__ == '__main__':
    app.run(debug=True)