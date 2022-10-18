// 将数值使用逗号隔开，一般用于金额的输入
export function getCommaNumber(value) {
  const list = value.toString().split('.');
  const prefix = list[0].charAt(0) === '-' ? '-' : '';
  let num = prefix ? list[0].slice(1) : list[0];
  let result = '';
  while (num.length > 3) {
    result = `,${num.slice(-3)}${result}`;
    num = num.slice(0, num.length - 3);
  }
  if (num) {
    result = num + result;
  }
  const listSecond = list[1] ? '.' + list[1] : '';
  return `${prefix}${result}${listSecond}`;
}


// 将金额转换为中文大写
export const moneyToCN = (n = 0) => {
  const fraction = ["角", "分"];
  const digit = ["零", "壹", "贰", "叁", "肆", "伍", "陆", "柒", "捌", "玖"];
  const unit = [
    ["元", "万", "亿"],
    ["", "拾", "佰", "仟"],
  ];
  const head = n < 0 ? "负" : "";
  n = Math.abs(n);
  let s = "";
  // 循环角、分，按照乘10%10的步骤对数值进行判断；
  for (let i = 0; i < fraction.length; i++) {
    s += (digit[Math.floor(n * 10 * Math.pow(10, i)) % 10] + fraction[i]).replace(/零./, "");
  }
  // 整数部分则角和分位无值，替换为整
  s = s || "整";
  n = Math.floor(n);
  // 整数部分%10余n,便在digit里循环n次，得到对应数字n的大写汉字，如果轮询到零就replace将零剪切，最后将得到的小数部位和整数部位进行拼接.
  for (let i = 0; i < unit[0].length && n > 0; i++) {
    let p = "";
    for (let j = 0; j < unit[1].length && n > 0; j++) {
      p = digit[n % 10] + unit[1][j] + p;
      n = Math.floor(n / 10);
    }
    s = p.replace(/(零.)*零$/, "").replace(/^$/, "零") + unit[0][i] + s;
  }
  return (
    head +
    s
    .replace(/(零.)*零元/, "元")
    .replace(/(零.)+/g, "零")
    .replace(/^整$/, "零元整")
  );
};
