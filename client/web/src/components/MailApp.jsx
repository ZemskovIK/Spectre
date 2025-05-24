import React, { useState, useEffect } from "react";
import axios from "axios";
import ModalWindow from "./ModalWindow";
import EmailList from "./EmailList";
import SendFormSection from "./SendFormSection";
import EditFormSection from "./EditFormSection";
import Auth from "./Auth";
import { jwtDecode } from "jwt-decode";

export default function MailApp() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [loadingAuth, setLoadingAuth] = useState(true);

  const [messages, setMessages] = useState([]);
  const [newMessage, setNewMessage] = useState("");
  const [author, setAuthor] = useState("");
  const [foundIn, setFoundIn] = useState("");
  const [dataInt, setDataInt] = useState("");
  const [newEditingMessage, setEditingNewMessage] = useState("");
  const [editingAuthor, setEditingAuthor] = useState("");
  const [editingFoundIn, setEditingFoundIn] = useState("");
  const [editingDataInt, setEditingDataInt] = useState("");
  const [loading, setLoading] = useState(false);
  const [selectedMessage, setSelectedMessage] = useState(null);
  const [showModal, setShowModal] = useState(false);
  const [deletingId, setDeletingId] = useState(null);
  const [isEditing, setIsEditing] = useState(false);
  const [editingId, setEditingId] = useState(null);
  const [isError, setIsError] = useState("Писем нет");
  const [searchQuery, setSearchQuery] = useState("");
  const [isAdmin, setIsAdmin] = useState(10);

  useEffect(() => {
    const checkAuth = async () => {
      const token = localStorage.getItem("token");
      if (!token) {
        setLoadingAuth(false);
        return;
      }

      try {
        const response = await axios.get("http://localhost:5000/api/letters", {
          headers: { Authorization: `Bearer ${token}` },
        });
        console.log(response);
        setIsAuthenticated(true);
        fetchMessages();
      } catch (err) {
        localStorage.removeItem("token");
      } finally {
        setLoadingAuth(false);
      }
    };

    checkAuth();
  }, []);

  useEffect(() => {
    const requestInterceptor = axios.interceptors.request.use((config) => {
      const token = localStorage.getItem("token");
      if (token) {
        const decoded = jwtDecode(token);
        setIsAdmin(decoded.role);
      }
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    });

    return () => {
      axios.interceptors.request.eject(requestInterceptor);
    };
  }, []);

  const handleLogout = () => {
    localStorage.removeItem("token");
    setIsAuthenticated(false);
    setMessages([]);
  };

  const fetchMessages = async () => {
    try {
      const response = await axios.get("http://localhost:5000/api/letters");
      if (response.data.error == null) {
        setMessages(response.data.content);
      } else {
        setIsError(response.data.error);
      }
    } catch (error) {
      console.error("Ошибка загрузки писем:", error);
      if (error.response && error.response.status === 401) {
        handleLogout();
      }
    }
    console.log(localStorage.getItem("token"));
    const token = localStorage.getItem("token");
    const secretKey = "test_secret";
    const decoded = jwtDecode(token);

    console.log(decoded.role);
    console.log();
  };

  const filteredMessages = searchQuery.trim()
    ? messages.filter((message) =>
        message?.author?.toLowerCase()?.includes(searchQuery.toLowerCase())
      )
    : messages;

  function transDataToDima(data) {
    let newVal = `${data.slice(6, 10)}-${data.slice(3, 5)}-${data.slice(
      0,
      2
    )}T00:00:00Z`;
    return newVal;
  }

  function handleChangeModes() {
    setIsEditing(!isEditing);
    setEditingDataInt("");
    setDataInt("");
  }

  const handleChange = (e) => {
    let value = e.target.value;
    value = value.replace(/[^\d]/g, "");

    if (value.length > 2) {
      value = `${value.slice(0, 2)}.${value.slice(2)}`;
    }
    if (value.length > 5) {
      value = `${value.slice(0, 5)}.${value.slice(5, 9)}`;
    }

    if (value.length > 10) return;

    setDataInt(value);
    setEditingDataInt(value);
  };

  const deleteMessage = async (id) => {
    setDeletingId(id);
    try {
      await axios.delete(`http://localhost:5000/api/letters/${id}`);
      fetchMessages();
    } catch (error) {
      console.error("Ошибка удаления:", error);
      if (error.response && error.response.status === 401) {
        handleLogout();
      }
    } finally {
      setMessages(messages.filter((msg) => msg.id !== id));
      setDeletingId(null);
    }
  };

  const editMessage = async (e) => {
    e.preventDefault();
    if (!newEditingMessage.trim() || !editingAuthor.trim()) return;

    setLoading(true);
    try {
      const response = await axios.put(
        `http://localhost:5000/api/letters/${editingId}`,
        {
          body: newEditingMessage,
          author: editingAuthor.trim(),
          found_at: transDataToDima(editingDataInt.trim()),
          found_in: editingFoundIn.trim(),
        }
      );
      fetchMessages();
      setMessages([...messages, response.data]);
      setEditingNewMessage("");
      setEditingAuthor("");
      setEditingDataInt("");
      setEditingFoundIn("");
      setEditingId("");
    } catch (error) {
      console.error("Ошибка редактирования:", error);
      if (error.response && error.response.status === 401) {
        handleLogout();
      }
    } finally {
      setLoading(false);
    }
  };

  const sendMessage = async (e) => {
    e.preventDefault();
    if (!newMessage.trim() || !author.trim()) return;

    setLoading(true);
    try {
      const response = await axios.post("http://localhost:5000/api/letters", {
        body: newMessage,
        author: author.trim(),
        found_at: transDataToDima(dataInt.trim()),
        found_in: foundIn.trim(),
      });
      fetchMessages();
      setMessages([...messages, response.data]);
      setNewMessage("");
      setAuthor("");
      setDataInt("");
      setFoundIn("");
    } catch (error) {
      console.error("Ошибка отправки:", error);
      if (error.response && error.response.status === 401) {
        handleLogout();
      }
    } finally {
      setLoading(false);
    }
  };

  const openMessage = (message) => {
    setSelectedMessage(message);
    setShowModal(true);
  };

  if (loadingAuth) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-xl">Проверка авторизации...</div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return (
      <Auth
        setIsAuthenticated={setIsAuthenticated}
        fetchMessages={fetchMessages}
      />
    );
  }

  return (
    <div className="min-h-screen m-6">
      {/* Кнопка выхода */}
      <div className="flex justify-end">
        <button
          onClick={handleLogout}
          className={`bg-red-500 mr-12 ${
            isAdmin != 6 && "mt-6"
          } text-white py-2 px-6 rounded-lg hover:bg-red-600 transition-colors absolute`}
        >
          Выйти
        </button>
      </div>

      <EmailList
        isAdmin={isAdmin == 6 ? true : false}
        searchQuery={searchQuery}
        setSearchQuery={setSearchQuery}
        isError={isError}
        messages={filteredMessages}
        onEmailClick={openMessage}
        onDeleteEmail={deleteMessage}
        deletingId={deletingId}
      ></EmailList>

      {isAdmin == 6 ? (
        <div className="lg:w-6xl 2xl:w-7xl mx-auto mt-6 bg-gradient-to-r from-gray-200 via-gray-350 to-gray-400 rounded-lg shadow-md p-6">
          <button
            onClick={handleChangeModes}
            className={`py-1 px-1 mr-3 rounded-lg text-white font-medium text-2xl font-semibold ${
              isEditing
                ? "bg-gray-500 hover:bg-gray-700"
                : "bg-gray-600 hover:bg-gray-700"
            }`}
          >
            Добавить новое письмо
          </button>{" "}
          <button
            onClick={handleChangeModes}
            className={`py-1 px-1 rounded-lg text-white font-medium text-2xl font-semibold mb-4 ${
              isEditing
                ? "bg-gray-600 hover:bg-gray-700"
                : "bg-gray-500 hover:bg-gray-700"
            }`}
          >
            Редактировать письмо
          </button>
          {!isEditing ? (
            <form
              onSubmit={sendMessage}
              className="space-y-2"
              autoComplete="off"
            >
              <SendFormSection
                handleChange={handleChange}
                author={author}
                setAuthor={setAuthor}
                foundIn={foundIn}
                setFoundIn={setFoundIn}
                dataInt={dataInt}
                setDataInt={setDataInt}
                newMessage={newMessage}
                setNewMessage={setNewMessage}
              ></SendFormSection>

              <button
                type="submit"
                disabled={
                  loading ||
                  !newMessage.trim() ||
                  !author.trim() ||
                  dataInt.length < 10
                }
                className={`px-6 py-2 rounded-lg text-white font-medium ${
                  loading ||
                  !newMessage.trim() ||
                  !author.trim() ||
                  dataInt.length < 10
                    ? "bg-gray-300 cursor-not-allowed"
                    : "bg-gray-600 hover:bg-gray-700"
                }`}
              >
                {loading ? "Отправка..." : "Отправить"}
              </button>
            </form>
          ) : (
            <form
              onSubmit={editMessage}
              className="space-y-2"
              autoComplete="off"
            >
              <EditFormSection
                handleChange={handleChange}
                editingId={editingId}
                setEditingId={setEditingId}
                editingAuthor={editingAuthor}
                editingDataInt={editingDataInt}
                editingFoundIn={editingFoundIn}
                newEditingMessage={newEditingMessage}
                setEditingAuthor={setEditingAuthor}
                setEditingDataInt={setEditingDataInt}
                setEditingFoundIn={setEditingFoundIn}
                setEditingNewMessage={setEditingNewMessage}
              ></EditFormSection>
              <button
                type="submit"
                disabled={
                  loading ||
                  !newEditingMessage.trim() ||
                  !editingAuthor.trim() ||
                  editingDataInt.length < 10
                }
                className={`px-6 py-2 rounded-lg text-white font-medium ${
                  loading ||
                  !newEditingMessage.trim() ||
                  !editingAuthor.trim() ||
                  editingDataInt.length < 10
                    ? "bg-gray-300 cursor-not-allowed"
                    : "bg-gray-600 hover:bg-gray-700"
                }`}
              >
                {loading ? "Отправка..." : "Сохранить"}
              </button>
            </form>
          )}
        </div>
      ) : (
        ""
      )}

      {showModal && (
        <ModalWindow
          onClose={() => setShowModal(false)}
          dataInt={selectedMessage.found_at}
          title={selectedMessage.author}
          foundIn={selectedMessage.found_in}
          text={selectedMessage.body}
          backgroundImage="https://static.vecteezy.com/system/resources/previews/032/048/239/non_2x/paper-vintage-background-recycle-brown-paper-crumpled-texture-ai-generated-free-photo.jpg"
        ></ModalWindow>
      )}
    </div>
  );
}
