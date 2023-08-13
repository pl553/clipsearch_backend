import numpy
import time
import csv
import torch
import clip
import os
import zmq
from PIL import Image
import io
import jsend

ZMQ_PORT = "5554"

device = "cuda" if torch.cuda.is_available() else "cpu"
model, preprocess = clip.load("ViT-L/14", device=device, download_root='models/')

THRESHOLD = 17

if os.path.isfile("image_features.pt"):
    image_features = torch.load('image_features.pt')
else:
    image_features = {}

context = zmq.Context()
socket = context.socket(zmq.REP)
socket.bind("tcp://*:" + ZMQ_PORT)

print("Image feature extraction daemon listening on tcp://localhost:" + ZMQ_PORT)

while True:
    message = socket.recv_multipart()
    try:
        image_id = message[0].decode("utf-8")
        image_bytes = message[1]

        image = Image.open(io.BytesIO(image_bytes))
        image = preprocess(image).unsqueeze(0).to(device)
        
        with torch.no_grad():
            # predict
            image_features[image_id] = model.encode_image(image)
            image_features[image_id] /= image_features[image_id].norm(dim=-1, keepdim=True)
        
        image_features[image_id] = image_features[image_id].cpu()

        torch.save(image_features, "image_features.pt")
    except Exception as e:
        socket.send_string(jsend.NewError(str(e)))
        continue
    socket.send_string(jsend.New(None))
