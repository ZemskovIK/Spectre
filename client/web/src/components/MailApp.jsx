import React, { useState, useEffect } from "react";
import axios from "axios";
import ModalWindow from "./ModalWindow";
import EmailList from "./EmailList";

export default function MailApp() {
  const [messages, setMessages] = useState([]);
  const [newMessage, setNewMessage] = useState("");
  const [author, setAuthor] = useState("");
  const [loading, setLoading] = useState(false);
  const [selectedMessage, setSelectedMessage] = useState(null);
  const [showModal, setShowModal] = useState(false);
  const [deletingId, setDeletingId] = useState(null);

  useEffect(() => {
    fetchMessages();
  }, []);

  const fetchMessages = async () => {
    try {
      const response = await axios.get("http://localhost:5000/messages");
      setMessages(response.data);
    } catch (error) {
      console.error("ошибка отправки:", error);
    }
  };

  const deleteMessage = async (id) => {
    setDeletingId(id);
    try {
      await axios.delete(`http://localhost:5000/messages/${id}`);
      setMessages(messages.filter((msg) => msg.id !== id));
    } catch (error) {
      console.error("ошибка удаления:", error);
    } finally {
      setDeletingId(null);
    }
  };

  const sendMessage = async (e) => {
    e.preventDefault();
    if (!newMessage.trim() || !author.trim()) return;

    setLoading(true);
    try {
      const response = await axios.post("http://localhost:5000/messages", {
        text: newMessage,
        author: author.trim(),
        // id2: "1",
      });
      setMessages([...messages, response.data]);
      setNewMessage("");
      setAuthor("");
    } catch (error) {
      console.error("ошибка отправки:", error);
    } finally {
      setLoading(false);
    }
  };

  const openMessage = (message) => {
    setSelectedMessage(message);
    setShowModal(true);
  };

  return (
    <div className="min-h-screen m-5">
      {/* <div className="max-w-6xl mx-auto mb-6">
        <h1 className="text-3xl font-bold text-blue-600">Spectre</h1>
        <p className="text-gray-600">какая нибудь шапка</p>
      </div> */}

      <EmailList
        messages={messages}
        onEmailClick={openMessage}
        onDeleteEmail={deleteMessage}
        deletingId={deletingId}
      ></EmailList>

      <div className="max-w-6xl mx-auto mt-6 bg-gradient-to-r from-gray-200 via-gray-350 to-gray-400 rounded-lg shadow-md p-6">
        <h2 className="text-xl font-semibold mb-4">Добавить новое письмо</h2>
        <form onSubmit={sendMessage} className="space-y-4" autocomplete="off">
          <div className="">
            <label
              htmlFor="author"
              className="block text-sm font-medium text-gray-700 mb-1"
            >
              Автор
            </label>
            <input
              id="author"
              type="text"
              value={author}
              onChange={(e) => setAuthor(e.target.value)}
              placeholder="автор"
              className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            />
          </div>
          <div>
            <label
              htmlFor="письмо"
              className="block text-sm font-medium text-gray-700 mb-1"
            >
              Письмо
            </label>
            <textarea
              id="message"
              value={newMessage}
              onChange={(e) => setNewMessage(e.target.value)}
              placeholder="пишите тут"
              className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              rows={4}
              required
            />
          </div>
          <button
            type="submit"
            disabled={loading || !newMessage.trim() || !author.trim()}
            className={`px-6 py-2 rounded-lg text-white font-medium ${
              loading || !newMessage.trim() || !author.trim()
                ? "bg-gray-300 cursor-not-allowed"
                : "bg-gray-600 hover:bg-gray-700"
            }`}
          >
            {loading ? "отправка..." : "отправить"}
          </button>
        </form>
      </div>

      {showModal && (
        <ModalWindow
          onClose={() => setShowModal(false)}
          title={selectedMessage.author}
          text={selectedMessage.text}
          backgroundImage="https://static.vecteezy.com/system/resources/previews/032/048/239/non_2x/paper-vintage-background-recycle-brown-paper-crumpled-texture-ai-generated-free-photo.jpg"
        ></ModalWindow>
      )}
    </div>
  );
}
