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
model, preprocess = clip.load("ViT-L/14@336px", device=device, download_root='models/')
    
context = zmq.Context()
socket = context.socket(zmq.REP)
socket.bind("tcp://*:" + ZMQ_PORT)

print("Image embedding daemon listening on tcp://localhost:" + ZMQ_PORT)

while True:
    try:
        image_bytes = socket.recv()

        image = Image.open(io.BytesIO(image_bytes))
        image = preprocess(image).unsqueeze(0).to(device)
        
        with torch.no_grad():
            image_embedding = model.encode_image(image)
            image_embedding /= image_embedding.norm(dim=-1, keepdim=True)
        
        socket.send_string(jsend.New(image_embedding.squeeze().tolist()))
    except Exception as e:
        socket.send_string(jsend.NewError(str(e)))
        continue
