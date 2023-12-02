import requests
from bs4 import BeautifulSoup
import os
import json

def isStringNull(string):
    return string != ''

os.mkdir('code_snippets')

file_path = 'questions_url.txt'

with open(file_path, 'r') as file:
    questions = file.readlines()

count = 1

for question in questions:

    res = requests.get(question.strip())

    soup = BeautifulSoup(res.text, "html.parser")

    title = soup.select_one(".question-hyperlink").text

    description = soup.select_one(".js-post-body").text.split("\n")

    description = list(filter(isStringNull, description))

    answer = soup.select(".js-post-body")[1].text.split("\n")

    answer = list(filter(isStringNull, answer))

    data = {
        "title": title,
        "description": description,
        "answer": answer
    }

    json_file_path = './code_snippets/q' + str(count) + '.json'

    with open(json_file_path, 'w') as json_file:
        json.dump(data, json_file, indent=4)

    count += 1
