#!/usr/bin/env python

from app import app


if __name__ == '__main__':
    app.run(port=app.config['WEB_PORT'], use_reloader=False)
