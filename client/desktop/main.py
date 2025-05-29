from  tkinter import *
from tkinter import ttk, messagebox, PhotoImage
import requests
import jwt
import json
import base64
import crypto
import os
from pyi_resource import resource_path

class MilitaryLettersApp:
    def __init__(self, root):
        self.root = root
        self.root.title("Военные письма")
        window_width = 800
        window_height = 600
        screen_width = self.root.winfo_screenwidth()
        screen_height = self.root.winfo_screenheight()
        x = (screen_width - window_width) // 2
        y = (screen_height - window_height) // 2
        self.root.geometry(f"{window_width}x{window_height}+{x}+{y}")
        self.img = PhotoImage(file=resource_path("letters.png"))
        self.root.iconphoto(True, self.img)
        
        self.style = ttk.Style()
        self.style.configure('TLabel', font=('Arial', 10))
        self.style.configure('TButton', font=('Arial', 10))
        self.style.configure('Header.TLabel', font=('Arial', 12, 'bold'))
        
        self.api_url = "http://localhost:5000"
        self.token = None
        self.current_user = None
        self.aes_key = b'\xb9M\x0b8\x00\x10\x90\x16\xc7\xed\x93\x08\xc1\x00J\xf2\x08\xb0\x01~\xb5_G\x805\xac\x95\xa2t`1\xde'
        self.hmac_key = b'Dp\xc2\xc6B\x16\xb8\\\xaf_z5\x8dC\x1f3\x19\n\xe1u8\xe1Q:\xd1}\xb2\xa0\xf8$\xa6\x0e'
        # self.aes_key = None
        # self.hmac_key = None
        self.client = None

        self.show_login_form()    
        
    def show_login_form(self):
        for widget in self.root.winfo_children():
            widget.destroy()
        
        login_frame = ttk.Frame(self.root, padding="20")
        login_frame.pack(fill=BOTH, expand=True)
        
        ttk.Label(login_frame, text="Логин:", style='Header.TLabel').grid(row=0, column=0, pady=5, sticky=W)
        self.login_entry = ttk.Entry(login_frame, width=30)
        self.login_entry.grid(row=0, column=1, pady=5, padx=5)
        
        ttk.Label(login_frame, text="Пароль:", style='Header.TLabel').grid(row=1, column=0, pady=5, sticky=W)
        self.password_entry = ttk.Entry(login_frame, width=30, show="*")
        self.password_entry.grid(row=1, column=1, pady=5, padx=5)
        
        login_btn = ttk.Button(login_frame, text="Войти", command=self.perform_login)
        login_btn.grid(row=2, column=1, pady=10, sticky=E)
    
    def perform_login(self):
        login = self.login_entry.get()
        password = self.password_entry.get()
        
        if not login or not password:
            error_msg = "Введите логин и пароль"
            messagebox.showerror("Ошибка", error_msg)
            return
        
        try:
            response = requests.post(
                f"{self.api_url}/login",
                json={"login": login, "password": password}
            )
            
            if response.status_code == 200:
                self.token = response.json().get("token")
                self.current_user = login
                
                if not self.token:
                    error_msg = "Токен отсутствует в ответе сервера"
                    messagebox.showerror("Ошибка", error_msg)
                    return
                
                try:
                    unverified_payload = jwt.decode(
                        self.token,
                        options={"verify_signature": False},
                        algorithms=["HS256"]
                    )
                    
                    if "sub" in unverified_payload:
                        if not isinstance(unverified_payload["sub"], str):
                            unverified_payload["sub"] = str(unverified_payload["sub"])
                    
                    payload = jwt.decode(
                        self.token,
                        "test_secret",
                        algorithms=["HS256"],
                        options={"verify_sub": False}
                    )
                    
                    self.user_role = payload.get("role", 1)
                    
                except Exception as e:
                    error_msg = f"Ошибка декодирования токена: {str(e)}"
                    self.user_role = 1
                    
                self.show_main_interface()
            else:
                error_msg = response.json().get("message", "Неверный логин или пароль")
                messagebox.showerror("Ошибка", error_msg)
                
        except requests.exceptions.RequestException as e:
            error_msg = f"Ошибка подключения: {str(e)}"
            messagebox.showerror("Ошибка", error_msg)
        except Exception as e:
            error_msg = f"Неожиданная ошибка: {str(e)}"
            messagebox.showerror("Ошибка", error_msg)
    
    def show_main_interface(self):
        for widget in self.root.winfo_children():
            widget.destroy()
        
        self.create_widgets()        

    def create_widgets(self):
        self.main_frame = ttk.Frame(self.root)
        self.main_frame.pack(fill=BOTH, expand=True, padx=10, pady=10)
        
        self.header_frame = ttk.Frame(self.main_frame)
        self.header_frame.pack(fill=X, pady=[0, 10])
        
        ttk.Label(self.header_frame, text="Архив", style='Header.TLabel').pack(side=LEFT)
        
        user_frame = ttk.Frame(self.header_frame)
        user_frame.pack(side=RIGHT)
        
        if self.current_user:
            ttk.Label(user_frame, text=self.current_user).pack(side=LEFT, padx=5)
        
        logout_btn = ttk.Button(user_frame, text="Выйти", command=self.logout)
        logout_btn.pack(side=RIGHT)
        
        self.tab_control = ttk.Notebook(self.main_frame)
        self.tab_control.pack(fill=BOTH, expand=True)
        
        self.read_tab = ttk.Frame(self.tab_control)
        self.tab_control.add(self.read_tab, text='Найти письмо')
        self.create_read_tab()
        
        if self.user_role == 6:
            self.create_tab = ttk.Frame(self.tab_control)
            self.tab_control.add(self.create_tab, text='Добавить письмо')
            self.create_add_tab()
            
            self.update_tab = ttk.Frame(self.tab_control)
            self.tab_control.add(self.update_tab, text='Изменить письмо')
            self.create_update_tab()
            
            self.delete_tab = ttk.Frame(self.tab_control)
            self.tab_control.add(self.delete_tab, text='Удалить письмо')
            self.create_delete_tab()
        
    def logout(self):
        self.token = None
        self.show_login_form()
    
    def make_authenticated_request(self, method, url, **kwargs):
        headers = kwargs.get('headers', {})
        headers['Authorization'] = f'Bearer {self.token}'
        kwargs['headers'] = headers
        
        try:
            response = requests.request(method, url, **kwargs)
            
            if response.status_code == 401:
                messagebox.showerror("Ошибка", "Сессия истекла. Пожалуйста, войдите снова.")
                self.logout()
                return None
            
            return response
        except requests.exceptions.RequestException as e:
            messagebox.showerror("Ошибка", f"Не удалось подключиться к серверу: {str(e)}")
            return None
    
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
        
        ttk.Label(frame, text="ID письма или автора:").grid(row=0, column=0, sticky=W, pady=5)
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
        
        print(f"main.py | submit_create content: {content}, {type(content)}")

        content_bytes = json.dumps(content).encode('utf-8')

        content_bytes = {
            "content": [base64.b64encode(content_bytes).decode('utf-8')]
        }

        print(f"main.py | submit_create content_bytes: {content_bytes}, {type(content_bytes)}")

        letter_data = encrypt(content_bytes, self.aes_key, self.hmac_key)
        
        try:
            response = self.make_authenticated_request(
                "POST", 
                f"{self.api_url}/api/letters",
                json=letter_data
            )
            response_data = response.json()


           
            if response.status_code == 200 and response_data.get("error") is None:
                messagebox.showinfo("Успех", f"Письмо успешно добавлено")
                self.clear_create_form()
            else:
                error_msg = response_data.get("error", "Неизвестная ошибка")
                messagebox.showerror("Ошибка", f"Ошибка: {error_msg}")
        except requests.exceptions.RequestException as e:
            messagebox.showerror("Ошибка", f"Не удалось подключиться к серверу: {str(e)}")
            
    def search_letter(self):  # GET /api/letters/{letter_id}
        query = self.read_query.get().strip()
        if not query:
            messagebox.showerror("Ошибка", "Введите ID письма или автора")
            return
        
        try:
            if query.isdigit():
                print(f"[DEBU] searchID")
                self._search_by_letter_id(query)
            else:
                print(f"[DEBU] searchAuthor")
                self._search_by_author(query)
        except Exception as e:
            messagebox.showerror("Ошибка", f"Произошла ошибка: {str(e)}")

    def _search_by_letter_id(self, letter_id):
        response = self.make_authenticated_request(
            "GET", 
            f"{self.api_url}/api/letters/{letter_id}"
        )
        response_data = decrypt(response.json(), self.aes_key, self.hmac_key)
        
        if response.status_code != 200:
            error_msg = response.json().get("error", "Не удалось получить список писем")
            messagebox.showerror("Ошибка", error_msg)
            return
        
        bytes_data = base64.b64decode(response_data.get("content")[0])
        json_data = json.loads(bytes_data.decode('utf-8'))
        
        if json_data:
            self.display_results([json_data])
        else:
            messagebox.showinfo("Информация", "Письмо с данным ID не найдено")

    def _search_by_author(self, author_name):
        response = self.make_authenticated_request(
            "GET", 
            f"{self.api_url}/api/letters"
        )
        response_data = decrypt(response.json(), self.aes_key, self.hmac_key)
        
        if response.status_code != 200:
            error_msg = response.json().get("error", "Не удалось получить список писем")
            messagebox.showerror("Ошибка", error_msg)
            return
        
        all_letters = []
        for item in response_data.get("content"):
                bytes_data = base64.b64decode(item)
                json_data = json.loads(bytes_data.decode('utf-8'))
                all_letters.append(json_data)
                                
        results = [
            letter for letter in all_letters 
            if str(letter.get("author", "")).lower() == author_name.lower()
        ]
        
        if results:
            self.display_results(results)
        else:
            messagebox.showinfo("Информация", "Письма данного автора не найдены")

    def get_all_letters(self): # GET /api/letters
        try:
            response = self.make_authenticated_request(
                "GET", 
                f"{self.api_url}/api/letters"
            )
            response_data = decrypt(response.json(), self.aes_key, self.hmac_key)

            result = []
            for item in response_data.get("content"):
                bytes_data = base64.b64decode(item)
                json_data = json.loads(bytes_data.decode('utf-8'))
                result.append(json_data)
                            
            self.results_body.config(state=NORMAL)
            self.results_body.delete("1.0", END)
            
            if response.status_code == 200 and response_data.get("content"):                
                if result:
                    for letter in result:
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
            # else:
            #     error_msg = letters.get("error", "Неизвестная ошибка")
            #     self.results_body.insert(END, f"Ошибка: {error_msg}")
                
            self.results_body.config(state=DISABLED)
        except requests.exceptions.RequestException as e:
            messagebox.showerror("Ошибка", f"Не удалось подключиться к серверу: {str(e)}")
 
    def fetch_letter_body(self): # GET /api/letters/{letter_id}
        letter_id = self.update_id.get()
        if not letter_id:
            messagebox.showerror("Ошибка", "Введите ID письма")
            return
            
        try:
            response = self.make_authenticated_request(
                "GET", 
                f"{self.api_url}/api/letters/{letter_id}"
            )
            response_data = decrypt(response.json(), self.aes_key, self.hmac_key)
            if response.status_code == 200:
                bytes_data = base64.b64decode(response_data.get("content")[0])
                letter_data = json.loads(bytes_data.decode('utf-8'))

                if letter_data:    
                    self.update_author.delete(0, END)
                    self.update_body.delete("1.0", END)
                    self.update_found_at.delete(0, END)
                    self.update_found_in.delete(0, END)
                    
                    self.update_author.insert(0, letter_data.get('author', ''))
                    self.update_body.insert("1.0", letter_data.get('body', ''))
                    self.update_found_at.insert(0, letter_data.get('found_at', '')[:10])
                    self.update_found_in.insert(0, letter_data.get('found_in', ''))
                else:
                    messagebox.showerror("Ошибка", "Данные письма не получены")
            else:
                messagebox.showinfo("Информация", "Письмо с данным ID не найдено")
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
        
        # Я(АЗАМАТ) ТРОГАЛ ЭТУ ЧАСТЬ КОДА!!
        content_bytes = json.dumps(content).encode('utf-8')

        content_bytes = {
            "content": [base64.b64encode(content_bytes).decode('utf-8')]
        }

        print(f"main.py | submit_update() content_bytes: {content_bytes}, {type(content_bytes)}")

        content = encrypt(content_bytes, self.aes_key, self.hmac_key)
        # Трогал только то, что между этим и верхним комментами.
        # Что я изменил?? Добавил шифрование в эту функцию.

        try:
            response = self.make_authenticated_request(
                "PUT", 
                f"{self.api_url}/api/letters/{letter_id}",
                json=content
            )
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
            response = self.make_authenticated_request(
                "DELETE", 
                f"{self.api_url}/api/letters/{letter_id}"
            )
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
    
def decrypt(data, aes_key, hmac_key):
    print(f"\nmain.py | decrypt() data: {data}, {type(data)}\n")

    crypto_box = crypto.Aes256CbcHmac(aes_key, hmac_key)

    decrypted_text = crypto_box.decrypt(data)

    data_list = json.loads(decrypted_text.decode('utf-8'))
    content_base64_list = [base64.b64encode(item.encode('utf-8')).decode('utf-8')
        for item in data_list]
    result = {
        "content": content_base64_list
    }

    print(f"\nmain.py | decrypt() result: {result}\n")
    return result

def encrypt(data, aes_key, hmac_key):
    # data = request.get_json()
    content = [base64.b64decode(data['content'][0]).decode("utf-8")]

    json_str = json.dumps(content)
    data = json_str.encode('utf-8')

    print(f"\nmain.py | encrypt() data in bytes: {data}, {type(data)}\n")
    # print(f"\nserver.py | encrypt() content: {content}\n")

    crypto_box = crypto.Aes256CbcHmac(aes_key, hmac_key)
    nonce = os.urandom(12)

    encrypted_text = crypto_box.encrypt(data, nonce)

    print(f"\nmain.py | encrypt() encrypted_text: {encrypted_text}, {type(encrypted_text)}\n")

    # data_list = json.loads(encrypted_text.decode('utf-8'))
    # content_base64_list = [base64.b64encode(item.encode('utf-8')).decode('utf-8')
    #     for item in data_list]
    # result = {
    #     "content": content_base64_list
    # }

    return json.dumps(encrypted_text)

def ecdh(server_pub):

    clinet = crypto.ECDHKeyExchange() # 4
    client_pub = client.get_public_key_base64() # 5

    client.compute_shared_secret(server_pub) # 7,9

    # Ключи снизу используем для шифрования и проверки целостности
    aes_key = clinet.aes_key
    hmac_key = client.hmac_key

    result = {
        "content": client_pub
    }

    print(f"main.py | ecdh() result: {result}, {type(result)}")
    
    return result

if __name__ == "__main__":
    root = Tk()
    app = MilitaryLettersApp(root)
    root.mainloop()