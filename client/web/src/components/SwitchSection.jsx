import React from "react";
import { Switch } from "antd";

export default function SwitchSection({ showContent, setShowContent }) {
  const onChange = (checked) => {
    if (showContent == "1") {
      setShowContent("2");
    } else {
      setShowContent("1");
    }
  };
  return <Switch defaultChecked onChange={onChange} />;
}
