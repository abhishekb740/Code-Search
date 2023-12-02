import requests
from bs4 import BeautifulSoup

for i in range(1, 2001):
    res = requests.get("https://stackoverflow.com/questions/tagged/java?sort=MostVotes&edited=true?page=" + str(i))

    soup  = BeautifulSoup(res.text, "html.parser")

    questions = soup.select(".s-post-summary")

    for question in questions:
        first_href = question.find('a')['href'] if questions and questions[0].find('a') else None
        print('https://stackoverflow.com' + first_href)
    