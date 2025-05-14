import tkinter as tk
from tkinter import ttk, messagebox, scrolledtext
import requests

API_URL = "http://localhost:8080/messages"

def clear_fields():
    entry_id.delete(0, tk.END)
    entry_title.delete(0, tk.END)
    entry_content.delete("1.0", tk.END)

def create_message():
    try:
        data = {
            "id": int(entry_id.get()),
            "title": entry_title.get(),
            "content": entry_content.get("1.0", tk.END).strip()
        }
        r = requests.post(API_URL, json=data)
        output.insert(tk.END, f"[CREATE] {r.status_code}: {r.text}\n")
    except Exception as e:
        output.insert(tk.END, f"❌ Ошибка: {e}\n")

def read_message():
    try:
        message_id = entry_id.get()
        r = requests.get(f"{API_URL}/{message_id}")
        if r.status_code == 200:
            msg = r.json()
            entry_title.delete(0, tk.END)
            entry_title.insert(0, msg['title'])
            entry_content.delete("1.0", tk.END)
            entry_content.insert("1.0", msg['content'])
            output.insert(tk.END, f"[READ] {r.status_code}: {msg}\n")
        else:
            output.insert(tk.END, f"[READ] {r.status_code}: {r.text}\n")
    except Exception as e:
        output.insert(tk.END, f"❌ Ошибка: {e}\n")

def update_message():
    try:
        message_id = entry_id.get()
        data = {
            "title": entry_title.get(),
            "content": entry_content.get("1.0", tk.END).strip()
        }
        r = requests.put(f"{API_URL}/{message_id}", json=data)
        output.insert(tk.END, f"[UPDATE] {r.status_code}: {r.text}\n")
    except Exception as e:
        output.insert(tk.END, f"❌ Ошибка: {e}\n")

def delete_message():
    try:
        message_id = entry_id.get()
        r = requests.delete(f"{API_URL}/{message_id}")
        output.insert(tk.END, f"[DELETE] {r.status_code}: {r.text}\n")
    except Exception as e:
        output.insert(tk.END, f"❌ Ошибка: {e}\n")

root = tk.Tk()
root.title("AmClient-PC")
root.geometry("700x500")

frm = ttk.Frame(root, padding=10)
frm.pack(fill="both", expand=True)

ttk.Label(frm, text="ID:").grid(column=0, row=0, sticky="e")
entry_id = ttk.Entry(frm, width=10)
entry_id.grid(column=1, row=0, sticky="w")

ttk.Label(frm, text="Название:").grid(column=0, row=1, sticky="e")
entry_title = ttk.Entry(frm, width=40)
entry_title.grid(column=1, row=1, sticky="w")

ttk.Label(frm, text="Содержимое:").grid(column=0, row=2, sticky="ne")
entry_content = tk.Text(frm, height=6, width=50)
entry_content.grid(column=1, row=2, sticky="w")

btn_frame = ttk.Frame(frm)
btn_frame.grid(column=1, row=3, pady=10)

ttk.Button(btn_frame, text="Создать", command=create_message).grid(column=0, row=0, padx=5)
ttk.Button(btn_frame, text="Прочитать", command=read_message).grid(column=1, row=0, padx=5)
ttk.Button(btn_frame, text="Обновить", command=update_message).grid(column=2, row=0, padx=5)
ttk.Button(btn_frame, text="Удалить", command=delete_message).grid(column=3, row=0, padx=5)

ttk.Button(btn_frame, text="Очистить поля", command=clear_fields).grid(column=4, row=0, padx=5)

output = scrolledtext.ScrolledText(frm, height=12)
output.grid(column=0, row=4, columnspan=2, pady=10)

root.mainloop()