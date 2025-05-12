import React, { useState, useEffect } from "react";
import axios from "axios";
import MessagesSection from "./MessagesSection";

export default function Chat() {
  const [messages, setMessages] = useState([]);
  const [newMessage, setNewMessage] = useState("");
  const [loading, setLoading] = useState(false);

  function isMobileDevice() {
    return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(
      navigator.userAgent
    );
  }

  useEffect(() => {
    fetchMessages();
  }, []);

  const fetchMessages = async () => {
    try {
      const response = await axios.get("http://localhost:5000/messages");
      setMessages(response.data);
    } catch (error) {
      // alert("не удалось прогрузить сообщения");
      console.error("Error fetching messages:", error);
    }
  };

  const sendMessage = async (e) => {
    e.preventDefault();
    if (!newMessage.trim()) return;

    setLoading(true);
    try {
      const response = await axios.post("http://localhost:5000/messages", {
        text: newMessage,
      });
      setMessages([...messages, response.data]);
      setNewMessage("");
    } catch (error) {
      alert("не удалось отправить сообщение");
      console.error("Error sending message:", error);
    } finally {
      setLoading(false);
    }
  };

  // return (
  //   <div className="bg-gradient-to-r from-gray-400 via-gray-450 to-gray-500">
  //     <div className="pt-5 flex-col h-screen">
  //       <div className="max-w-4xl mx-auto bg-white rounded-lg shadow-md overflow-hidden">
  //         <div className="p-4 bg-gray-600 text-white">
  //           <h1 className="text-2xl font-bold">Спектр</h1>
  //         </div>

  //         <div className="">
  //           <div className="p-4 h-140 overflow-y-auto">
  //             {messages.length === 0 ? (
  //               <p className="text-gray-500 text-center py-8">сообщений нет</p>
  //             ) : (
  //               <ul className="space-y-3">
  //                 {messages.map((message) => (
  //                   <MessagesSection
  //                     id={message.id}
  //                     timestamp={message.timestamp}
  //                     text={message.text}
  //                   ></MessagesSection>
  //                 ))}
  //               </ul>
  //             )}
  //           </div>
  //           <form onSubmit={sendMessage} className="p-4 border-t">
  //             <div className="flex space-x-2">
  //               <input
  //                 type="text"
  //                 value={newMessage}
  //                 onChange={(e) => setNewMessage(e.target.value)}
  //                 placeholder="писать"
  //                 className="flex-1 px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
  //                 disabled={loading}
  //               />
  //               <button
  //                 type="submit"
  //                 disabled={loading || !newMessage.trim()}
  //                 className={`px-4 py-2 rounded-lg text-white ${
  //                   loading || !newMessage.trim()
  //                     ? "bg-gray-300 cursor-not-allowed"
  //                     : "bg-gray-600 hover:bg-gray-700"
  //                 }`}
  //               >
  //                 {loading ? "отправка..." : "отправить"}
  //               </button>
  //             </div>
  //           </form>
  //         </div>
  //       </div>
  //     </div>
  //   </div>
  // );

  return (
    <div className="">
      <div
        className={`flex justify-center flex-col h-screen bg-gradient-to-r from-gray-400 via-gray-450 to-gray-500
          ${isMobileDevice() ? "p-2 pl-2 pr-2" : "p-4 pl-50 pr-50"}`}
      >
        <div className=" p-5 bg-gray-600 text-white flex justify-center rounded-t-xl">
          <h1 className="text-2xl  font-bold">Спектр</h1>
        </div>

        <div className="flex-1 overflow-hidden flex flex-col">
          <div className="flex-1 overflow-y-auto p-4 bg-gradient-to-r from-gray-200 via-gray-250 to-gray-300">
            {messages.length === 0 ? (
              <div className="h-full flex items-center justify-center">
                <p className="text-gray-500">сооьщений нет</p>
              </div>
            ) : (
              <ul className="space-y-3">
                //{" "}
                {messages.map((message) => (
                  <MessagesSection
                    id={message.id}
                    timestamp={message.timestamp}
                    text={message.text}
                  ></MessagesSection>
                ))}
              </ul>
            )}
          </div>

          <form
            onSubmit={sendMessage}
            className={`p-4 border-t bg-white rounded-b-xl ${
              isMobileDevice() && "mb-15"
            }`}
          >
            <div className="flex space-x-2">
              <input
                type="text"
                value={newMessage}
                onChange={(e) => setNewMessage(e.target.value)}
                placeholder="писать"
                className="flex-1 px-4 border rounded-lg focus:outline-none focus:ring-2 focus:ring-gray-500"
                disabled={loading}
              />
              <button
                type="submit"
                disabled={loading || !newMessage.trim()}
                className={`px-4 py-2 rounded-lg text-white ${
                  loading || !newMessage.trim()
                    ? "bg-gray-300 cursor-not-allowed"
                    : "bg-gray-600 hover:bg-gray-700"
                }`}
              >
                {loading ? "отправка..." : "отправить"}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}
