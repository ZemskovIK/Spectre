import React from "react";
import { useState } from "react";
export default function EmailList({
  isError,
  messages,
  onEmailClick,
  onDeleteEmail,
  deletingId,
}) {
  const chunkArray = (arr, size) => {
    return Array.from({ length: Math.ceil(arr.length / size) }, (v, i) =>
      arr.slice(i * size, i * size + size)
    );
  };

  const messageRows = chunkArray(messages, 4);

  return (
    <div className="lg:w-6xl 2xl:w-7xl mx-auto bg-gradient-to-r from-gray-200 via-gray-350 to-gray-400 rounded-lg shadow-md overflow-hidden">
      <div className="p-2 border-b">
        <h2 className="text-lg font-semibold">Письма ({messages.length})</h2>
      </div>
      <div className="overflow-y-auto" style={{ height: "300px" }}>
        <div className="p-4 space-y-4">
          {isError != "Писем нет" || messages.length === 0 ? (
            <div className="text-gray-500 py-8 text-center">{isError}</div>
          ) : (
            messageRows.map((row, rowIndex) => (
              <div key={rowIndex} className="grid grid-cols-4 gap-4">
                {row.map((message) => (
                  <div
                    key={message.id}
                    className="relative group relative bg-gray-450 p-4 rounded-lg shadow-sm border cursor-pointer hover:shadow-md transition-transform hover:scale-105"
                    onClick={() => onEmailClick(message)}
                  >
                    <button
                      onClick={(e) => {
                        e.stopPropagation();
                        onDeleteEmail(message.id);
                      }}
                      className=" rounded-full absolute right-0 text-red-400 p-1 group-hover:opacity-100 transition-opacity hover:bg-red-600 z-10"
                      disabled={deletingId === message.id}
                    >
                      <strong className="opacity-0 group-hover:opacity-100 transition-opacity duration-200 ">
                        &#10060;
                      </strong>
                    </button>
                    <div className="flex flex-col items-center">
                      <img src="https://e7.pngegg.com/pngimages/178/760/png-clipart-paper-envelope-letter-mail-envelope-miscellaneous-material.png"></img>
                      <span className="font-medium text-sm text-center truncate w-full">
                        {message.author}
                      </span>
                      <span className="text-xs text-gray-500 mt-1">
                        номер: {message.id}
                      </span>
                    </div>
                  </div>
                ))}
                {row.length < 4 &&
                  Array(4 - row.length)
                    .fill(0)
                    .map((_, i) => (
                      <div
                        key={`empty-${i}`}
                        className="opacity-0"
                        aria-hidden="true"
                      >
                        <div className="p-4 rounded-lg border">
                          <div className="h-10 w-10 mb-2"></div>
                        </div>
                      </div>
                    ))}
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  );
}
