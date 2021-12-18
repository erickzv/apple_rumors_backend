from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

import apple_news

app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_methods=["OPTIONS", "POST", "GET", "DELETE"],
    allow_headers=["*"],
)

news_urls = [
    apple_news.mac_rumors,
    apple_news.apple_insiders,
    apple_news.nine_to_five_mac
]


@app.get("/all_news")
def all_news():
    domains, news = [], []
    for data in news_urls:
        news.append(apple_news.get_news(**data))
        domains.append(data["url"][8:-4])
    news.append(domains)
    return news


@app.get("/macrumors")
def mac_rumors():
    news = apple_news.get_news(**apple_news.mac_rumors)
    return news


@app.get("/apple_insiders")
def apple_insiders():
    news = apple_news.get_news(**apple_news.apple_insiders)
    return news


@app.get("/9to5mac")
def nine_to_five_mac():
    news = apple_news.get_news(**apple_news.nine_to_five_mac)
    return news
