import asyncio
from typing import Dict

import aiohttp
from bs4 import BeautifulSoup

mac_rumors = {
    "url": "https://macrumors.com",
    "find_tag": "h2"
}

apple_insiders = {
    "url": "https://appleinsider.com",
    "find_tag": "h2"
}

nine_to_five_mac = {
    "url": "https://9to5mac.com",
    "find_tag": "h1"
}


async def parse_news(session: aiohttp.ClientSession, url: str, find_tag: str):
    async with session.get(url) as response:
        soup = BeautifulSoup(await response.text(), "lxml")
        tags = soup.find_all(find_tag)

        rumors = []
        for tag in tags:
            if (a_tag := tag.find("a")) is not None:
                rumors.append({"title": a_tag.text, "href": a_tag.get("href")})
            if len(rumors) == 10:
                break
        return rumors


async def get_news(news_urls: Dict[str, str]):
    tasks, domains = [], []
    async with aiohttp.ClientSession() as session:
        for data in news_urls:
            tasks.append(parse_news(session, **data))
            domains.append(data["url"][8:-4])
        return await asyncio.gather(*tasks), domains
