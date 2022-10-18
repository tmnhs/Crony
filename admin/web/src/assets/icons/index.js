const context = require.context('./modules', false, /\.svg$/)
context.keys().map(context)
