// 计算浏览器滚动条宽度
let scrollbarWidth;

export default function getScrollbarWidth() {
  if (scrollbarWidth !== undefined) return scrollbarWidth;
  const div = document.createElement('div');
  div.style.width = '100px';
  div.style.height = '100px';
  div.style.position = 'absolute';
  div.style.top = '-9999px';
  div.style.overflow = 'scroll';
  document.body.appendChild(div);
  scrollbarWidth = (div.offsetWidth - div.clientWidth);
  document.body.removeChild(div);
  return scrollbarWidth;
}