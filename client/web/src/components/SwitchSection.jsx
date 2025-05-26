import React from "react";
import { Switch, ConfigProvider } from "antd";

export default function SwitchSection({
  showContent,
  setShowContent,
  onClick,
}) {
  const onChange = (checked) => {
    if (showContent == "1") {
      setShowContent("2");
    } else {
      setShowContent("1");
    }
    console.log(showContent);
  };
  return (
    <ConfigProvider
      theme={{
        components: {
          Switch: {
            colorPrimary: "#a2a4a9",
            colorPrimaryHover: "#000000",
            colorHover: "#000000",
          },
        },
      }}
    >
      <Switch defaultChecked onChange={onChange} />
    </ConfigProvider>
  );
}
