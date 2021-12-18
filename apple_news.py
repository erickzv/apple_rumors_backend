import requests
from bs4 import BeautifulSoup
from pydantic import BaseModel

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


class News(BaseModel):
    title: str
    href: str


def get_news(url: str, find_tag: str):
    try:
        response = requests.get(url)
    except requests.exceptions.RequestException:
        return []

    soup = BeautifulSoup(response.text, "html.parser")
    h2_tags = soup.find_all(find_tag)

    rumors = []
    for tag in h2_tags:
        if (a_tag := tag.find("a")) is not None:
            rumors.append(News(title=a_tag.text, href=a_tag.get("href")))
        if len(rumors) == 10:
            break
    return rumors
