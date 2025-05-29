import sys, os

def resource_path(relative_path, folder=None):
    try:
        base_path = sys._MEIPASS
    except Exception:
        base_path = os.path.abspath(".")
    if folder:
        return os.path.join(base_path, folder, relative_path)
    else:
        return os.path.join(base_path, relative_path)