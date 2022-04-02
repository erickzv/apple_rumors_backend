from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import ORJSONResponse

import apple_news

app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_methods=["OPTIONS", "POST", "GET", "DELETE"],
    allow_headers=["*"],
)

news_urls = (
    apple_news.mac_rumors,
    apple_news.apple_insiders,
    apple_news.nine_to_five_mac
)


@app.get("/all_news", response_class=ORJSONResponse)
async def all_news():
    news, domains = await apple_news.get_news(news_urls)
    return {
        "news": news,
        "websites": domains
    }
