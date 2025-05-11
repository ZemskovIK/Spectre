import React, { useState, useEffect } from "react";
import axios from "axios";
import MessagesSection from "./components/MessagesSection";
import Chat from "./components/Chat";
import MailApp from "./components/MailApp";

export default function App() {
  // const imageUrl =
  //   "https://tranceam.org/wp-content/uploads/2021/04/WW2-header.jpg";

  return (
    // <Chat></Chat>

    // <div
    //   style={{ "--image-url": `url(${imageUrl})` }}
    //   className="min-h-screen bg-[image:var(--image-url)] bg-cover bg-center"
    // >
    <MailApp></MailApp>
    // </div>
  );
}
