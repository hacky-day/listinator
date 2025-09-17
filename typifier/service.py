from fastapi import FastAPI
from pydantic import BaseModel
import joblib

# load the model
model = joblib.load("model.pkl")

app = FastAPI()

TYPE_UUID_MAP = {
    "miscellaneous": "c29ebd85-812e-4cf6-bfc7-c8368eb83334",
    "fruit & vegetables & nuts": "fe0b085b-2df9-4422-a7cb-7867947719a5",
    "canned goods": "0c9b99fb-c2c8-41e4-8afa-b8cca3ac2ca1",
    "sauces & spices & dressings": "7c693d05-4939-44e6-845d-57951720e886",
    "drinks & alcohol": "0828b46f-98c9-41ea-9918-164751782861",
    "bakery": "e693272f-4a40-4c0e-9e38-8ebb33004271",
    "spreads": "ab8328c2-29e2-4767-a6fb-27d8e11dc8df",
    "coffee & tea": "21b7a2d6-0507-41dc-9a41-4f8a3c86564a",
    "cereals & muesli": "97ef6e7e-6c1a-47bc-9d34-35e26a1a0d5c",
    "pasta & rice": "1a78a64a-ff86-49db-b64d-45a8b2e76c25",
    "cooking & baking": "5d6b6b67-34f3-4a48-bb63-cf65f0f2219d",
    "meat & fish": "d67bd9ce-56f1-4227-885b-0656f74edb22",
    "frozen": "b98f7846-a4cd-4b00-86bf-a6714e982469",
    "dairy & chilled": "36298b3b-fcd5-4189-b34f-dae3dea08412",
    "sweets & snacks": "13f6bd3e-aeeb-4890-955f-fd91c2450a7e",
    "household & baby & pets": "a14bca10-13b7-4a9c-a663-75a5203c3f09",
    "ready meals & broth & sauce": "3de4b7ac-60be-4d65-8dbf-431f2c6d1270",
}


class Item(BaseModel):
    product: str


@app.post("/")
def predict(item: Item):
    produkt_lower = item.product.lower()
    type = model.predict([produkt_lower])[0]
    uuid = TYPE_UUID_MAP.get(type, None)
    return {"product": item.product, "type": type, "uuid": uuid}
