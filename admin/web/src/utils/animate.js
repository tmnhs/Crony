
/**
 * 常见的动画算法
 * @param {Number} t 已经经过的时间
 * @param {Number} b 起始位置
 * @param {Number} c 总距离
 * @param { Number} d 运动总共要用的时间
 * @param { Number} s 缓动的强度，数值越大效果越明显
 */
// 总距离:总时间 = 移动了多少距离:经过了多少时间
const tween = {
  // 匀速运动
  linear(t, b, c, d) {
    return c * t / d + b;
  },
  // 先慢后快
  easeIn(t, b, c, d, s = 4) {
    return c * (t /= d) * Math.pow(t, s) + b;
  },
  // 先快后慢
  easeOut(t, b, c, d, s = 4) {
    return c * ((t = t / d - 1) * Math.pow(t, s) + 1) + b;
  },
  // 前半段先慢后快，后半段先快后慢
  easeInOut(t, b, c, d) {
    let time = t;
    time /= d / 2;
    if (time < 1) {
      return (c / 2) * time * time + b;
    }
    time--;
    return (-c / 2) * ((time - 2) * time - 1) + b;
  }
};

//定义Animate类
const Animate = function (elem) {
  this.elem = elem; // 进行运动的 dom 节点
  this.startTime = 0; // 动画开始时间
  this.startPos = 0; // 动画开始时，dom 节点的位置，即 dom 的初始位置
  this.endPos = 0; // 动画结束时，dom 节点的位置，即 dom 的目标位置
  this.propertyName = ''; // dom 节点需要被改变的 css 属性名
  this.easing = () => { }; // 缓动算法
  this.duration = 0; // 动画持续时间
  this.callback = () => { }; //动画结束后的回调
};

//启动动画
Animate.prototype.start = function (propertyName, endPos, duration, easing, callback) {
  this.startTime = Date.now();
  this.startPos = parseInt(window.getComputedStyle(this.elem)[propertyName]);
  this.propertyName = propertyName;
  this.endPos = endPos;
  this.duration = duration;
  this.easing = tween[easing];
  if (callback && typeof callback === 'function') {
    this.callback = callback;
  }
  this.step();
};

//每次运动要做的事情
Animate.prototype.step = function () {
  const passTime = Date.now() - this.startTime;
  const nextPosition = this.easing(passTime, this.startPos, this.endPos - this.startPos, this.duration);
  this.elem.style[this.propertyName] = nextPosition + 'px';
  if (passTime < this.duration) {
    window.requestAnimationFrame(this.step.bind(this));
  } else {
    this.elem.style[this.propertyName] = this.endPos + 'px';
    this.callback();
  }
};

export default Animate;
