import pandas as pd
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.linear_model import LogisticRegression
from sklearn.pipeline import make_pipeline
import joblib

# read data
data = pd.read_csv("training.csv")

# create pipeline
model = make_pipeline(
    TfidfVectorizer(analyzer="char_wb", ngram_range=(2, 4)),
    LogisticRegression(max_iter=200),
)

# train
model.fit(data["product"], data["type"])

# store model
joblib.dump(model, "model.pkl")
print("model saved")

# test
test_produkte = [
    "butter",
    "margarine",
    "baguette",
    "apfel",
    "shampoo",
    "tomaten",
    "ketchup",
    "cola",
    "müsli",
    "spaghetti",
    "zucker",
    "rindfleisch",
    "blumenkohl",
    "schokolade",
    "kaffee",
    "babyflasche",
    "fertigsuppe",
]

for t in test_produkte:
    predicted = model.predict([t])[0]
    print(f"{t} -> Type: {predicted}")
