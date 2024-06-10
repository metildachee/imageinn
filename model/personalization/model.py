from flask import Flask, request, jsonify
from surprise import Dataset, Reader
from surprise import SVD
import pandas as pd

app = Flask(__name__)

# Load the ratings dataset
ratings_file_path = '../../data/movielens/ratings.csv'
reader = Reader(line_format='user item rating timestamp', sep=',', rating_scale=(1, 5), skip_lines=1)
data = Dataset.load_from_file(ratings_file_path, reader=reader)
trainset = data.build_full_trainset()

# Load the movies dataset
movies_file_path = '../../data/movielens/movies.csv'
movies = pd.read_csv(movies_file_path)

# Initialize the SVD algorithm
algo = SVD()

@app.route('/predict/<int:user_id>', methods=['GET'])
def predict(user_id):
    try:
        # Predict ratings for all movies
        all_movies = movies['movieId'].unique()
        predictions = [(movie_id, algo.predict(trainset.to_inner_uid(user_id), movie_id).est) for movie_id in all_movies]

        # Convert ratings to binary (1 for liked, 0 for not liked)
        liked_movies = [movie_id for movie_id, rating in predictions if rating >= 4]

        # Get movie titles for the liked movies
        liked_movie_titles = movies[movies['movieId'].isin(liked_movies)]['title']

        return jsonify({"user_id": user_id, "liked_movies": liked_movie_titles.tolist()})
    except ValueError as e:
        return jsonify({"error": str(e)}), 400

@app.route('/train_model/<int:user_id>', methods=['POST'])
def train_model(user_id):
    try:
        # Get clicked items from request
        clicked_items = request.json.get('clicked_items')
        if not clicked_items:
            return jsonify({"error": "No clicked items provided"}), 400

        # Add new user ratings to the trainset
        for item_id in clicked_items:
            trainset.add_user(trainset.to_inner_uid(user_id), trainset.to_inner_iid(item_id), 1.0)  # Assuming binary rating

        # Retrain the model
        algo.fit(trainset)

        return jsonify({"message": "Model retrained successfully"})
    except ValueError as e:
        return jsonify({"error": str(e)}), 400

if __name__ == '__main__':
    app.run(debug=True)
