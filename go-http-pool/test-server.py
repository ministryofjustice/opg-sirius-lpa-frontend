from flask import Flask, request
app = Flask(__name__)

@app.route("/")
def hello_world():
    return {"a": "booch", "b": 2}

@app.route("/goodbye_cruel_world")
def goodbye_cruel_world():
    return {"c": 10, "d": 20}

@app.route("/hello_sunshine")
def hello_sunshine():
    return {"e": 10, "f": {"g": "777", "h": "333"}}

@app.route("/say_it")
def say_it():
    to_say = request.args.get("say_this")
    return {"say": to_say}

if __name__ == '__main__':
    app.run()
