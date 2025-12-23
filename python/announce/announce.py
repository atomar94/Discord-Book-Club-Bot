from schedule.schedule import read_schedule

from schedule.schedule import get_next_book, update_book, update_schedule, get_book

def make_announcement(dryrun=True):
    """Checks every day if it is Monday. If so, posts the schedule."""
    # Our Annoucements channel ID
    ANNOUNCEMENT_CHANNEL_ID = 1451785614163709964 
    
    # 1. Get the structured schedule data
    full_schedule = read_schedule()
    
    if not full_schedule:
        return ""

    # if uninitialized, assign a book now.
    upcoming = full_schedule[0]
    if not upcoming.get("book_name", ""):
        book = get_next_book()
        book["read"] = True
        if not dryrun:
            update_book(book)
        upcoming["book_name"] = book["name"]
        if not dryrun:
            update_schedule(upcoming)

    # if we are 2 weeks from finishing this book then pick another
    plus_2 = full_schedule[2]
    if not plus_2.get("book_name", ""):
        book = get_next_book()
        book["read"] = True
        if not dryrun:
            update_book(book)
        plus_2["book_name"] = book["name"]
        if not dryrun:
            for s in full_schedule[2:6]:
                s["book_name"] = book["name"]
                update_schedule(s)

    # 1. Get the structured schedule data
    full_schedule = read_schedule()
    
    if not full_schedule:
        return ""

    print(full_schedule)
    book = get_book(full_schedule[0]["book_name"])
    book_initials = book_shortname(book["name"])

    link_text = f" - [Link]({book['link']})" if book.get("link", "") else "" 

    # 2. Build the message (showing this week + next 4)
    response_lines = [
        "📢 **Weekly Book Club Schedule Update**",
        f"This month we are reading **{book['name']}**{link_text}",
        "",
        "There is no sign up needed. Please feel free to drop in!",
        "",
        "Here is where we are headed for the next month:",
        "---"
    ]

    for i, item in enumerate(full_schedule[:5]):
        prefix = "📍 **THIS WEEK:**" if i == 0 else f"{i}. "
        line = f"{prefix} `{item['date']}` — **{item['cafe']}** - [Directions]({item['url']}) - Reading: {book_initials}"
        response_lines.append(line)


    final_message = "\n".join(response_lines)
    return final_message

def book_shortname(book_name):
    tokens = book_name.split(' ')
    initials = []
    LOWER_CASE_INITIALISM_WORDS = ["of", "the", "a", "an", "but", "and", "or", "in",
                                   "at", "yet", "for", "to"]
    for token in tokens:
        if token.lower() in LOWER_CASE_INITIALISM_WORDS:
            initials.append(token[0].lower())
        else:
            initials.append(token[0].upper())

    return ".".join(initials[:4]) + "."