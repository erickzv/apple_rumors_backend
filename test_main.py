from fastapi.testclient import TestClient

from main import app

client = TestClient(app)

def test_all_news():
  response = client.get("/all_news")
  assert response.status_code == 200
  assert type(response.json()) == list
