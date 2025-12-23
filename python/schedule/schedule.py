import json
import os
import re
import random
import logging
from datetime import datetime, timedelta

def make_schedule_object(name, url, book_name=""):
    return {
        "name": name,
        "url": url,
        "book_name": book_name,
    }

def push_schedule(name, url):
    schedule_table = None
    with open("schedule.json", "r") as schedule_file:
        try:
            schedule_table = json.loads(schedule_file.read())

        except Exception as e:
            logging.error(e.with_traceback())
            raise e
        
    cafe_pool = schedule_table["cafe_pool"]
    schedule = schedule_table["schedule"]

    pool_names = set(entry["name"].lower().strip() for entry in cafe_pool)
    if name.lower().strip() not in pool_names:
        cafe_pool.append(make_schedule_object(name, url))

    schedule_names = set(entry["name"].lower().strip() for entry in schedule)
    if name.lower().strip() not in schedule_names:
        schedule.append(make_schedule_object(name, url))

    # "w" is overwrite
    with(open("schedule.json", "w")) as schedule_file:
        schedule_file.write(json.dumps(schedule_table))


def read_schedule():
    """
    Returns a list of dictionaries containing the cafe name and its meeting date.
    Each item is scheduled for a consecutive Saturday.
    """
    schedule_table = None
    with open("schedule.json", "r") as schedule_file:
        try:
            schedule_table = json.loads(schedule_file.read())

        except Exception as e:
            logging.error(e.with_traceback())
            raise e

    schedule = schedule_table["schedule"]
    
    if not schedule:
        return []

    # 1. Find the next Saturday from today
    today = datetime.now()
    # weekday() returns 0 for Monday ... 5 for Saturday, 6 for Sunday
    # (5 - today.weekday()) % 7 gives the days until next Saturday
    days_until_saturday = (5 - today.weekday()) % 7
    
    # If today is Saturday, we might want to stay on today or skip to next week
    # If it's already past meeting time today, we add 7 days.
    if days_until_saturday == 0:
        # Optional logic: skip to next week if it's currently late Saturday
        pass 

    base_date = today + timedelta(days=days_until_saturday)

    # 2. Iterate through the queue and increment the week for each cafe
    schedule_list = []
    for index, entry in enumerate(schedule):
        meeting_date = base_date + timedelta(weeks=index)
        
        schedule_list.append({
            "cafe": entry["name"],
            "url": entry["url"],
            "date": meeting_date.strftime("%A, %B %d, %Y"),
            "book_name": entry["book_name"]
        })

    return schedule_list


def update_schedule(schedule):
    schedule_table = None
    with open("schedule.json", "r") as schedule_file:
        try:
            schedule_table = json.loads(schedule_file.read())

        except Exception as e:
            logging.error(e.with_traceback())
            raise e
        
    logging.info("about to access schedule_table")
    full_schedule = schedule_table["schedule"]

    for entry in full_schedule:
        if schedule["date"].lower().strip() == entry["date"].lower().strip():
            entry = schedule

    # "w" is overwrite
    with(open("schedule.json", "w")) as schedule_file:
        print(schedule_table)
        schedule_file.write(json.dumps(schedule_table))

def make_book_object(name, link=None, description=None, votes=0):
    return {
        "name": name,
        "link": link or "No link provided",
        "description": description or "No description provided",
        "votes": votes,
        "read": False
    }


def update_book(book):
    schedule_table = None
    try:
        with open("schedule.json", "r") as f:
            schedule_table = json.load(f)
    except Exception as e:
        logging.error(f"Error reading file: {e}")
        return

    # Initialize book_pool if it doesn't exist in an old version of the file
    if "book_pool" not in schedule_table:
        schedule_table["book_pool"] = []

    book_pool = schedule_table["book_pool"]
    for entry in book_pool:
        if book["name"].lower().strip() == entry["name"].lower().strip():
            entry = book # update

    # Save back to file
    with open("schedule.json", "w") as f:
        json.dump(schedule_table, f, indent=4)

def push_book(name, link=None, description=None):
    """Logs a new book recommendation into the book_pool."""
    schedule_table = None
    try:
        with open("schedule.json", "r") as f:
            schedule_table = json.load(f)
    except Exception as e:
        logging.error(f"Error reading file: {e}")
        return

    # Initialize book_pool if it doesn't exist in an old version of the file
    if "book_pool" not in schedule_table:
        schedule_table["book_pool"] = []

    book_pool = schedule_table["book_pool"]

    # Check for duplicates (case-insensitive)
    book_names = set(entry["name"].lower().strip() for entry in book_pool)
    if name.lower().strip() not in book_names:
        new_book = make_book_object(name, link, description)
        book_pool.append(new_book)

    # Save back to file
    with open("schedule.json", "w") as f:
        json.dump(schedule_table, f, indent=4)
    
def update_book_votes(book_name, votes):
    """
    Finds a book in the book_pool and updates its vote count.
    """
    logging.info(f"update_book_values({book_name}, {votes})")
    schedule_table = None
    try:
        with open("schedule.json", "r") as f:
            schedule_table = json.load(f)
    except Exception as e:
        logging.error(f"Error reading file for vote update: {e}")
        return

    if "book_pool" not in schedule_table:
        logging.warning("No book_pool found in schedule.json")
        return

    # Look for the book (case-insensitive and stripped of whitespace)
    book_found = False
    target_name = book_name.lower().strip()

    for book in schedule_table["book_pool"]:
        if book["name"].lower().strip() == target_name:
            book["votes"] = votes
            book_found = True
            break

    if book_found:
        with open("schedule.json", "w") as f:
            json.dump(schedule_table, f, indent=4)
        logging.info(f"Updated votes for '{book_name}' to {votes}")
    else:
        logging.warning(f"Book '{book_name}' not found in the pool.")

def get_book(book_name):
    schedule_table = None
    try:
        with open("schedule.json", "r") as f:
            schedule_table = json.load(f)
    except Exception as e:
        logging.error(f"Error reading file for top books: {e}")
        return []

    book_pool = schedule_table.get("book_pool", [])

    if not book_pool:
        return None

    for book in book_pool:
        if book["name"].lower().strip() == book_name:
            return book

    return None
   

def get_next_book():
    """
    Returns the next book to read.
    """
    schedule_table = None
    try:
        with open("schedule.json", "r") as f:
            schedule_table = json.load(f)
    except Exception as e:
        logging.error(f"Error reading file for top books: {e}")
        return []

    book_pool = schedule_table.get("book_pool", [])

    if not book_pool:
        return []

    # Sort the list of dictionaries by the 'votes' key
    # reverse=True ensures the highest votes come first
    book_pool = list(filter(lambda x: not x.get('read', False), book_pool))
    book_pool = sorted(book_pool, key=lambda x: x['votes'], reverse=True)

    return book_pool[random.randint(0,4)]
