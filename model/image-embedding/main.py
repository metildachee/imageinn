from flask import Flask, request, jsonify
import torch
from transformers import CLIPProcessor, CLIPModel
from PIL import Image
import torchvision.transforms as transforms

app = Flask(__name__)

device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
model_name = "openai/clip-vit-base-patch32"
clip = CLIPModel.from_pretrained(model_name).to(device)
processor = CLIPProcessor.from_pretrained(model_name)

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

if __name__ == '__main__':
    app.run(debug=True)