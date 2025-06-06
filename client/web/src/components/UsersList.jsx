import React from "react";
import { useState } from "react";
export default function UsersList({
  isAdmin,
  searchQuery,
  setSearchQuery,
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
    <>
      <div className="pl-4 pr-8 pt-2">
        <input
          type="text"
          placeholder="Поиск по логину..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="w-full p-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-gray-400"
        />
      </div>
      <div
        className="overflow-y-auto"
        style={isAdmin ? { height: "260px" } : { height: "560px" }}
      >
        <div className="p-4 space-y-4">
          {isError != "Писем нет" || messages.length === 0 ? (
            <div className="text-gray-500 py-8 text-center">{isError}</div>
          ) : (
            messageRows.map((row, rowIndex) => (
              <div key={rowIndex} className="grid grid-cols-4 gap-4">
                {row.map((message) => (
                  <div
                    key={message.id}
                    className="relative group relative bg-gray-450 p-4 rounded-lg shadow-sm border cursor-pointer hover:shadow-md transition-transform h-57 hover:scale-105"
                    onClick={() => onEmailClick(message)}
                  >
                    {isAdmin && (
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
                    )}
                    <div className="flex flex-col items-center">
                      {/* https://e7.pngegg.com/pngimages/178/760/png-clipart-paper-envelope-letter-mail-envelope-miscellaneous-material.png */}
                      <img
                        src="https://cdn2.iconfinder.com/data/icons/users-22/64/_User_Man_Circle-1024.png"
                        className="z-0 group-hover:z-[-1] group-hover:-rotate-30 transition-transform duration-350"
                        style={{ height: "150px" }}
                      ></img>
                      <span className="font-medium text-sm text-center truncate w-full">
                        {message.login}
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
    </>
  );
}
