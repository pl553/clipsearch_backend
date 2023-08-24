import numpy
import time
import csv
import torch
import clip
import os
import zmq
import jsend

ZMQ_PORT = "5553"

device = "cuda" if torch.cuda.is_available() else "cpu"
model, preprocess = clip.load("ViT-L/14@336px", device=device, download_root='models/')

THRESHOLD = 17

context = zmq.Context()
socket = context.socket(zmq.REP)
socket.bind("tcp://*:" + ZMQ_PORT)

print("Text embedding daemon listening on tcp://localhost:" + ZMQ_PORT)

while True:
    prompt = socket.recv().decode("utf-8")
    try:
        with torch.no_grad():
            text_features = model.encode_text(clip.tokenize(prompt))
            text_features /= text_features.norm(dim=-1, keepdim=True)
        socket.send_string(jsend.New(text_features.squeeze().tolist()))
    except Exception as e:
        socket.send_string(jsend.NewError(str(e)))
        continue
