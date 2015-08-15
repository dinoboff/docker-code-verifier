#!/usr/bin/env python3
#
import argparse
import base64
import json
import logging
import os
import subprocess
import sys

from flask import Flask, request, make_response

from codeverifier import TestRunner


JSON_TYPE = "application/json"
JS_TYPE = "application/javascript"
TIMEOUT = 5
RUNNER_SCRIPT = os.path.join(
    os.path.dirname(os.path.abspath(__file__)),
    'runner.py'
)


ERR_JSONREQUEST_REQUIRED = "jsonrequest is a required property."
ERR_JSONREQUEST_REQUIRED_BASE64 = (
    "jsonrequest is required property. It should be json object,"
    " base64 encoded or plain."
)
ERR_SOLUTION_REQUIRED = "The request didn't include a user solution."

CODE_OK = 200
CODE_BAD_REQUEST = 400

# Flask app
app = Flask(__name__)


# Comand line parser
parser = argparse.ArgumentParser(
    description=(
        "Sever running some python code and test it with some doctest tests."
    )
)

parser.add_argument("-d", "--debug", action='store_true')
parser.add_argument("-q", "--quiet", action='store_true')
parser.add_argument("-v", "--verbose", action='store_true')
parser.add_argument(
    "--host", default="localhost", help="Host to bind the server too"
)
parser.add_argument(
    "--port", type=int, default=5000, help="port to bind the server too"
)


def spawn(solution, tests=None):
    if tests is None:
        args = [sys.executable, RUNNER_SCRIPT, solution]
    else:
        args = [sys.executable, RUNNER_SCRIPT, solution, tests]

    proc = subprocess.Popen(
        args,
        stdout=subprocess.PIPE,
        universal_newlines=True,
    )

    try:
        out, _ = proc.communicate(timeout=TIMEOUT)
    except subprocess.TimeoutExpired:
        proc.kill()
        return {
            'solved': False,
            'errors': 'timeout'
        }

    if proc.returncode != 0:
        raise Exception('Failed to run tests.')

    return json.loads(out)


def response(code=CODE_OK, cb=None, **ctx):
    logging.debug(cb)
    output = "%s(%s)" % (cb, json.dumps(ctx),) if cb else json.dumps(ctx)
    resp = make_response(output, code)
    resp.headers["Content-Type"] = JS_TYPE if cb else JSON_TYPE
    resp.headers["Access-Control-Allow-Origin"] = "*"
    resp.status_code = code
    return resp


@app.route("/python", methods=["GET"])
@app.route("/python3", methods=["GET"])
def verifiy_get():
    req = request.args.get("jsonrequest")
    if req is None:
        return response(code=CODE_BAD_REQUEST, errors=ERR_JSONREQUEST_REQUIRED)

    try:
        req = base64.b64decode(req).decode("utf8")
    except Exception as e:
        logging.debug(e)

    try:
        req = json.loads(req)
    except Exception:
        return response(
            code=CODE_BAD_REQUEST,
            errors=ERR_JSONREQUEST_REQUIRED_BASE64
        )

    solution = req.get("solution")
    if not solution:
        return response(
            code=CODE_BAD_REQUEST,
            errors=ERR_SOLUTION_REQUIRED
        )

    resp = spawn(solution, req.get("tests"))
    return response(**resp)


@app.route("/python", methods=["POST"])
@app.route("/python3", methods=["POST"])
def verifiy_post():
    if request.headers.get("Content-Type").split(';')[0] == JSON_TYPE:
        req = request.get_json()
        if not req:
            return response(code=CODE_BAD_REQUEST, errors="Invalid JSON body.")
    else:
        req = request.form.get("jsonrequest")
        if req is None:
            return response(
                code=CODE_BAD_REQUEST, errors=ERR_JSONREQUEST_REQUIRED
            )

        try:
            req = json.loads(req)
        except Exception as e:
            logging.debug(e)
            return response(
                code=CODE_BAD_REQUEST, errors=ERR_JSONREQUEST_REQUIRED
            )

    solution = req.get("solution")
    if not solution:
        return response(
            code=CODE_BAD_REQUEST, errors=ERR_SOLUTION_REQUIRED
        )

    resp = spawn(solution, req.get("tests"))
    return response(**resp)


if __name__ == "__main__":
    args = parser.parse_args()
    if args.verbose:
        logging.basicConfig(level=logging.DEBUG)
    elif args.quiet:
        logging.basicConfig(level=logging.ERROR)
    else:
        logging.basicConfig(level=logging.INFO)
    app.run(host=args.host, port=args.port, debug=args.debug)
