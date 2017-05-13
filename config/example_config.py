"""
Configuration variables for application
Rename this file to config.py before use
"""

#Flask
DEBUG = True
APPLICATION_DIRECTORY = "/Users/autoturret/Documents/basecoat"
SECRET_KEY = "longrandomstringhere"

#Database
SQLALCHEMY_DATABASE_URI = 'postgresql://basecoat:mysupersecretdbpassword@postgres:5432/basecoat'
SQLALCHEMY_TRACK_MODIFICATIONS = False
