// 策略对象，封装校验规则
const strategies = {
  // 手机号
  mobile(value) {
    return /(^1[345789]\d{9}$)/.test(value);
  },
  // 电子邮箱
  email(value) {
    return /^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$/.test(value);
  },
  // 身份证号
  IDCard(value) {
    return /^(\d{6})(\d{4})(\d{2})(\d{2})(\d{3})([0-9]|X)$/.test(value);
  },
  // 最少长度
  min(value, limit) {
    return value.length >= limit;
  },
  // 最大长度
  max(value, limit) {
    return value.length <= limit;
  },
  // 是否必填
  required(value, limit) {
    if (!limit) {
      return true;
    };
    return value ? true : false;
  },
  // 自定义正则表达式校验
  pattern(value, limit) {
    return limit.test(value);
  }
}




// 存储校验器对象
let validators = [];


/**
 * 添加校验器
 * @param value{any}   要校验的值
 * @param rules{Array} 校验的规则。
 *
 */
function addValidator(value, rules) {
  rules.forEach(rule => {
    if (!rule.message) {
      throw new Error('缺少提示信息');
    }

    const message = rule.message;
    delete rule.message;
    let strategy = Object.keys(rule)[0];
    if (!strategy) {
      throw new Error('缺少校验类型');
    }

    if (strategy === 'type') {
      strategy = rule.type;
    }

    if (!Object.keys(strategies).includes(strategy)) {
      throw new Error(`校验器无法校验${strategy}类型`);
    }

    validators.push({
      strategy,
      value,
      message,
      limit: rule[strategy]
    })
  })
}


function validate(data, rules) {
  for (const [key, value] of Object.entries(rules)) {
    addValidator(data[key], value);
  }
  let errorMsg = '';
  validators.some(validator => {
    const { strategy, value, message, limit } = validator;
    const result = strategies[strategy](value, limit);
    if (!result) {
      errorMsg = message;
      return true;
    }
  })
  validators = [];
  return errorMsg;
}

export default validate;
