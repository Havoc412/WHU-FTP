// tag File Type
type IndexedType = {
  [key: string]: number;
};

export const fileType: IndexedType = {
  // info 需要和 fileicon.vue 组件数组联动
  dir: 0,
  txt: 1,
  doc: 2,
  "": 3,
};

export const fileTypeDescribe = ["文件夹", "文本", "DOC文档", "未定义"];

export const fileIconColor = [
  { unselected: "#999", selected: "#FFE391" },
  { unselected: "#999", selected: "#eeeeee" },
  { unselected: "#999", selected: "#A7F6F8" },
  { unselected: "#999", selected: "#eeeeee" },
];

// tag common hover color setting
export const iconHoverColor = {
  unselected: "#777",
  selected: "#fff",
};

// tag path recommand
export const localPath = [
  { path: "", name: "C:/" },
  { path: "C:/", name: "Users" },
  { path: "C:/Users/", name: "Desktop" },
  { path: "C:/Users/Havoc/Desktop/", name: "GoProject" },
];

export const ftpServerPath = [
  { path: "", name: "" },
  { path: "/", name: "TEST-1" },
];