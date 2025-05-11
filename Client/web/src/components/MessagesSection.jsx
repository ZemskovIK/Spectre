export default function MessagesSection({ id, timestamp, text }) {
  return (
    <>
      <li
        key={id}
        className="p-3 bg-gradient-to-r from-gray-300 to-gray-500 rounded-lg"
      >
        <div className="text-sm text-gray-500">
          {new Date(timestamp).toLocaleString()}
        </div>
        <div className="mt-1 break-words whitespace-pre-wrap overflow-x-auto">
          {text}
        </div>
      </li>
    </>
  );
}
