from  tkinter import *
from tkinter import ttk, messagebox
import requests

class MilitaryLettersApp:
    def __init__(self, root):
        self.root = root
        self.root.title("Военные письма")
        self.root.geometry("800x600+350+100")
        
        self.style = ttk.Style()
        self.style.configure('TLabel', font=('Arial', 10))
        self.style.configure('TButton', font=('Arial', 10))
        self.style.configure('Header.TLabel', font=('Arial', 12, 'bold'))
        
        self.api_url = "http://localhost:5000/api/letters"
        self.create_widgets()
    
    def create_widgets(self):
        self.main_frame = ttk.Frame(self.root)
        self.main_frame.pack(fill=BOTH, expand=True, padx=10, pady=10)
        
        self.header_frame = ttk.Frame(self.main_frame)
        self.header_frame.pack(fill=X, pady=[0, 10])
        
        ttk.Label(self.header_frame, text="Архив", style='Header.TLabel').pack(side=LEFT)
        
        self.tab_control = ttk.Notebook(self.main_frame)
        self.tab_control.pack(fill=BOTH, expand=True)
        
        self.create_tab = ttk.Frame(self.tab_control)
        self.tab_control.add(self.create_tab, text='Добавить письмо')
        self.create_add_tab()
        
        self.read_tab = ttk.Frame(self.tab_control)
        self.tab_control.add(self.read_tab, text='Найти письмо')
        self.create_read_tab()
        
        self.update_tab = ttk.Frame(self.tab_control)
        self.tab_control.add(self.update_tab, text='Изменить письмо')
        self.create_update_tab()
        
        self.delete_tab = ttk.Frame(self.tab_control)
        self.tab_control.add(self.delete_tab, text='Удалить письмо')
        self.create_delete_tab()
    
    def create_add_tab(self):
        frame = ttk.Frame(self.create_tab)
        frame.pack(fill=BOTH, expand=True, padx=10, pady=10)
        
        ttk.Label(frame, text="Автор:").grid(row=0, column=0, sticky=W, pady=5)
        self.create_author = ttk.Entry(frame, width=50)
        self.create_author.grid(row=0, column=1, pady=5, padx=5)
        
        ttk.Label(frame, text="Текст письма:").grid(row=1, column=0, sticky=W, pady=5)
        self.create_body = Text(frame, width=50, height=10)
        self.create_body.grid(row=1, column=1, pady=5, padx=5)
        
        ttk.Label(frame, text="Дата находки (ГГГГ-ММ-ДД):").grid(row=2, column=0, sticky=W, pady=5)
        self.create_found_at = ttk.Entry(frame, width=50)
        self.create_found_at.grid(row=2, column=1, pady=5, padx=5)
        
        ttk.Label(frame, text="Место находки:").grid(row=3, column=0, sticky=W, pady=5)
        self.create_found_in = ttk.Entry(frame, width=50)
        self.create_found_in.grid(row=3, column=1, pady=5, padx=5)
        
        submit_btn = ttk.Button(frame, text="Добавить письмо", command=self.submit_create)
        submit_btn.grid(row=4, column=1, pady=10, sticky=E)
    
    def create_read_tab(self):
        frame = ttk.Frame(self.read_tab)
        frame.pack(fill=BOTH, expand=True, padx=10, pady=10)
        
        ttk.Label(frame, text="ID письма:").grid(row=0, column=0, sticky=W, pady=5)
        self.read_query = ttk.Entry(frame, width=50)
        self.read_query.grid(row=0, column=1, pady=5, padx=5)
        
        search_btn = ttk.Button(frame, text="Найти", command=self.search_letter)
        search_btn.grid(row=0, column=2, padx=5)
        
        ttk.Label(frame, text="Результаты:").grid(row=1, column=0, sticky=W, pady=5)
        
        self.results_body = Text(frame, width=70, height=15, state=DISABLED)
        self.results_body.grid(row=2, column=0, columnspan=3, pady=5)
        
        get_all_btn = ttk.Button(frame, text="Показать все письма", command=self.get_all_letters)
        get_all_btn.grid(row=3, column=1, pady=10)
    
    def create_update_tab(self):
        frame = ttk.Frame(self.update_tab)
        frame.pack(fill=BOTH, expand=True, padx=10, pady=10)
        
        ttk.Label(frame, text="ID письма для обновления:").grid(row=0, column=0, sticky=W, pady=5)
        self.update_id = ttk.Entry(frame, width=50)
        self.update_id.grid(row=0, column=1, pady=5, padx=5)
        
        fetch_btn = ttk.Button(frame, text="Получить данные", command=self.fetch_letter_body)
        fetch_btn.grid(row=0, column=2, padx=5)
        
        ttk.Label(frame, text="Автор:").grid(row=1, column=0, sticky=W, pady=5)
        self.update_author = ttk.Entry(frame, width=50)
        self.update_author.grid(row=1, column=1, pady=5, padx=5)
        
        ttk.Label(frame, text="Текст письма:").grid(row=2, column=0, sticky=W, pady=5)
        self.update_body = Text(frame, width=50, height=10)
        self.update_body.grid(row=2, column=1, pady=5, padx=5)
        
        ttk.Label(frame, text="Дата находки (ГГГГ-ММ-ДД):").grid(row=3, column=0, sticky=W, pady=5)
        self.update_found_at = ttk.Entry(frame, width=50)
        self.update_found_at.grid(row=3, column=1, pady=5, padx=5)
        
        ttk.Label(frame, text="Место находки:").grid(row=4, column=0, sticky=W, pady=5)
        self.update_found_in = ttk.Entry(frame, width=50)
        self.update_found_in.grid(row=4, column=1, pady=5, padx=5)
        
        update_btn = ttk.Button(frame, text="Обновить письмо", command=self.submit_update)
        update_btn.grid(row=5, column=1, pady=10, sticky=E)
    
    def create_delete_tab(self):
        frame = ttk.Frame(self.delete_tab)
        frame.pack(fill=BOTH, expand=True, padx=10, pady=10)
        
        ttk.Label(frame, text="ID письма для удаления:").grid(row=0, column=0, sticky=W, pady=5)
        self.delete_id = ttk.Entry(frame, width=50)
        self.delete_id.grid(row=0, column=1, pady=5, padx=5)
        
        delete_btn = ttk.Button(frame, text="Удалить", command=self.delete_letter)
        delete_btn.grid(row=0, column=2, padx=5)
        
        ttk.Label(frame, text="Статус:").grid(row=1, column=0, sticky=W, pady=5)
        
        self.delete_status = Text(frame, width=70, height=5, state=DISABLED)
        self.delete_status.grid(row=2, column=0, columnspan=3, pady=5)
    
    def submit_create(self): # POST /api/letters
        content = {
            "author": self.create_author.get(),
            "body": self.create_body.get("1.0", END).strip(),
            "found_at": f"{self.create_found_at.get()}T00:00:00Z",
            "found_in": self.create_found_in.get()
        }
        
        if not all(content.values()):
            messagebox.showerror("Ошибка", "Все поля должны быть заполнены")
            return
        
        try:
            response = requests.post(self.api_url, json=content)
            response_data = response.json()
           
            if response.status_code == 200 and response_data.get("error") is None:
                messagebox.showinfo("Успех", f"Письмо успешно добавлено")
                self.clear_create_form()
            else:
                error_msg = response_data.get("error", "Неизвестная ошибка")
                messagebox.showerror("Ошибка", f"Ошибка: {error_msg}")
        except requests.exceptions.RequestException as e:
            messagebox.showerror("Ошибка", f"Не удалось подключиться к серверу: {str(e)}")
    
    def search_letter(self): # GET /api/letters/{letter_id}
        query = self.read_query.get()
        if not query:
            messagebox.showerror("Ошибка", "Введите ID письма")
            return
        
        try:
            if query.isdigit():
                response = requests.get(f"{self.api_url}/{query}")
                response_data = response.json()
                if response.status_code == 200 and response_data.get("content"):
                    self.display_results([response_data["content"]])
                else:
                    error_msg = response_data.get("error", "Письмо не найдено")
                    messagebox.showerror("Ошибка", error_msg)
                    self.display_results([])
        except requests.exceptions.RequestException as e:
            messagebox.showerror("Ошибка", f"Не удалось подключиться к серверу: {str(e)}")

    def get_all_letters(self): # GET /api/letters
        try:
            response = requests.get(f"{self.api_url}")
            response_data = response.json()
            
            self.results_body.config(state=NORMAL)
            self.results_body.delete("1.0", END)
            
            if response.status_code == 200 and response_data.get("content"):                
                letters = response_data["content"]
                if letters:
                    for letter in letters:
                        formatted_letter = (
                            f"ID: {letter.get('id', 'N/A')}\n"
                            f"Автор: {letter.get('author', 'N/A')}\n"
                            f"Дата: {letter.get('found_at', 'N/A')[:10]}\n"
                            f"Место: {letter.get('found_in', 'N/A')}\n"
                            f"Текст: {letter.get('body', 'N/A')}\n"
                            f"{'-'*30}\n"
                        )
                        self.results_body.insert(END, formatted_letter)
                else:
                    self.results_body.insert(END, "В базе нет писем")
            else:
                error_msg = letters.get("error", "Неизвестная ошибка")
                self.results_body.insert(END, f"Ошибка: {error_msg}")
                
            self.results_body.config(state=DISABLED)
        except requests.exceptions.RequestException as e:
            messagebox.showerror("Ошибка", f"Не удалось подключиться к серверу: {str(e)}")
 
    def fetch_letter_body(self): # GET /api/letters/{letter_id}
        letter_id = self.update_id.get()
        if not letter_id:
            messagebox.showerror("Ошибка", "Введите ID письма")
            return
            
        try:
            response = requests.get(f"{self.api_url}/{letter_id}")
            response_data = response.json()
            
            if response.status_code == 200 and response_data.get("content"):
                letter = response_data["content"]
                
                self.update_author.delete(0, END)
                self.update_body.delete("1.0", END)
                self.update_found_at.delete(0, END)
                self.update_found_in.delete(0, END)
                
                self.update_author.insert(0, letter.get('author', ''))
                self.update_body.insert("1.0", letter.get('body', ''))
                self.update_found_at.insert(0, letter.get('found_at', '')[:10])
                self.update_found_in.insert(0, letter.get('found_in', ''))
                
            else:
                error_msg = response_data.get("error", "Письмо не найдено")
                messagebox.showerror("Ошибка", error_msg)
        except requests.exceptions.RequestException as e:
            messagebox.showerror("Ошибка", f"Не удалось подключиться к серверу: {str(e)}")
    
    def submit_update(self): # PUT /api/letters/{letter_id}
        letter_id = self.update_id.get()
        if not letter_id:
            messagebox.showerror("Ошибка", "Введите ID письма")
            return
        
        content = {
            "author": self.update_author.get(),
            "body": self.update_body.get("1.0", END).strip(),
            "found_at": f"{self.update_found_at.get()}T00:00:00Z",
            "found_in": self.update_found_in.get()
        }
        
        if not all(content.values()):
            messagebox.showerror("Ошибка", "Все поля должны быть заполнены")
            return
        
        try:
            response = requests.put(f"{self.api_url}/{letter_id}", json=content)
            response_data = response.json()
            
            if response.status_code == 200 and response_data.get("error") is None:
                messagebox.showinfo("Успех", "Письмо успешно обновлено")
                self.clear_update_form()
            else:
                error_msg = response_data.get("error", "Неизвестная ошибка")
                messagebox.showerror("Ошибка", f"Ошибка: {error_msg}")
        except requests.exceptions.RequestException as e:
            messagebox.showerror("Ошибка", f"Не удалось подключиться к серверу: {str(e)}")
    
    def delete_letter(self): # DELETE /api/letters/{letter_id}
        letter_id = self.delete_id.get()
        if not letter_id:
            messagebox.showerror("Ошибка", "Введите ID письма")
            return
        
        try:
            response = requests.delete(f"{self.api_url}/{letter_id}")
            response_data = response.json()
            
            if response.status_code == 200 and response_data.get("error") is None:
                self.display_delete_status("Удаление успешно выполнено")
            else:
                error_msg = response_data.get("error", "Неизвестная ошибка")
                self.display_delete_status(f"Ошибка: {error_msg}")
        except requests.exceptions.RequestException as e:
            self.display_delete_status(None, f"Не удалось подключиться к серверу: {str(e)}")
    
    def display_results(self, letters):
        self.results_body.config(state=NORMAL)
        self.results_body.delete("1.0", END)
        
        if not letters:
            self.results_body.insert(END, "Письма не найдены")
        else:
            for letter in letters:
                formatted_letter = (
                    f"ID: {letter.get('id', 'N/A')}\n"
                    f"Автор: {letter.get('author', 'N/A')}\n"
                    f"Дата: {letter.get('found_at', 'N/A')[:10]}\n"
                    f"Место: {letter.get('found_in', 'N/A')}\n"
                    f"Текст: {letter.get('body', 'N/A')}\n"
                    f"{'-'*30}\n"
                )
                self.results_body.insert(END, formatted_letter)
        
        self.results_body.config(state=DISABLED)
    
    def display_delete_status(self, message):
        self.delete_status.config(state=NORMAL)
        self.delete_status.delete("1.0", END)
        self.delete_status.insert(END, message)
        self.delete_status.config(state=DISABLED)
    
    def clear_create_form(self):
        self.create_author.delete(0, END)
        self.create_body.delete("1.0", END)
        self.create_found_at.delete(0, END)
        self.create_found_in.delete(0, END)
        
    def clear_update_form(self):
        self.update_id.delete(0, END)
        self.update_author.delete(0, END)
        self.update_body.delete("1.0", END)
        self.update_found_at.delete(0, END)
        self.update_found_in.delete(0, END)
    

if __name__ == "__main__":
    root = Tk()
    app = MilitaryLettersApp(root)
    root.mainloop()