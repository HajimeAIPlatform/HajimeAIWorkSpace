
def remove_non_utf8(text):
    """
    Remove all non-UTF8 supported characters from a text string.

    Parameters:
    text (str): The input string containing UTF-8 and non-UTF8 characters.

    Returns:
    str: The cleaned string with only UTF-8 characters.
    """

    # Remove non-breaking spaces
    text = text.replace("\xa0", " ")
    # Encode to UTF-8 and ignore errors, then decode back
    cleaned_text = text.encode("utf-8", "ignore").decode("utf-8")
    return cleaned_text