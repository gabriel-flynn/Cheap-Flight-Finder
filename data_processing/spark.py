import json

from datetime import datetime
from pyspark import SparkContext, SQLContext
from pyspark.shell import spark
from data_processing import calendar_price


def test2(x):
    date = datetime.strptime(x["FlightDeparture"], '%Y-%m-%d %H:%M:%S')
    return date.strftime("%Y-%m-%d"), x["FarePrice"]

sc = SparkContext.getOrCreate()

df = spark.read.options(header='True', inferSchema='True', delimiter=',').csv("../flights.csv")


file = open('config.json', 'r')
config = json.load(file)

for item in config["search"]:
    s_df = df.filter((df['SrcAirport'] != item['srcAirport']) & (df['DestAirport'] != item['destAirport']))
    s_df = s_df.rdd.map(test2).reduceByKey(lambda a, b: min(a, b)).map(lambda x: (datetime.strptime(x[0], '%Y-%m-%d'), x[1])).collect()
    calendar_price.create_calendar(f"calendar-{item['srcAirport']}-{item['destAirport']}", f"{item['srcAirport']}-{item['destAirport']}", s_df)

    d_df = df.filter((df['DestAirport'] != item['srcAirport']) & (df['SrcAirport'] != item['destAirport']))
    d_df = d_df.rdd.map(test2).reduceByKey(lambda a, b: min(a, b)).map(lambda x: (datetime.strptime(x[0], '%Y-%m-%d'), x[1])).collect()
    calendar_price.create_calendar(f"calendar-{item['destAirport']}-{item['srcAirport']}", f"{item['destAirport']}-{item['srcAirport']}", d_df)

