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
model, preprocess = clip.load("ViT-L/14", device=device, download_root='models/')

THRESHOLD = 17

context = zmq.Context()
socket = context.socket(zmq.REP)
socket.bind("tcp://*:" + ZMQ_PORT)

print("Image search query daemon listening on tcp://localhost:" + ZMQ_PORT)

while True:
    prompt = socket.recv().decode("utf-8")
    resp_data = []
    try:
        if not os.path.isfile("image_features.pt"):
            socket.send_string(jsend.New([]))
            continue
        
        image_features = torch.load('image_features.pt')

        with torch.no_grad():
            zeroshot_weights = []
            texts = [prompt]
            texts = clip.tokenize(texts).cpu()
            class_embeddings = model.encode_text(texts) #embed with text encoder
            class_embeddings /= class_embeddings.norm(dim=-1, keepdim=True)
            class_embedding = class_embeddings.mean(dim=0)
            class_embedding /= class_embedding.norm()
            zeroshot_weights.append(class_embedding)
            zeroshot_weights = torch.stack(zeroshot_weights, dim=1).cpu()

        for image_id in image_features:
            logits = 100. * image_features[image_id].numpy() @ zeroshot_weights.numpy()
            resp_data.append({
                "id": int(image_id),
                "score": float(logits[0][0])
            })

        resp_data.sort(key=lambda image: image["score"], reverse=True)
    except Exception as e:
        socket.send_string(jsend.NewError(str(e)))
        continue
    socket.send_string(jsend.New({
        "totalCount": len(image_features),
        "images": resp_data
    }))
