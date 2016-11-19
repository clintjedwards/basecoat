"""
Configuration variables for application
Rename this file to config.py before use
"""

#Flask
DEBUG = True
APPLICATION_DIRECTORY = "/Users/autoturret/Documents/basecoat"

#Database
SQLALCHEMY_DATABASE_URI = 'sqlite:////' +  APPLICATION_DIRECTORY + '/database/basecoat.db'
SQLALCHEMY_TRACK_MODIFICATIONS = False
