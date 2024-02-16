export default (config = {}) => ({
    fields: {},
    hasErrors: false,
    pristine: true,
    valid: false,
    init() {
        for (const field in config) {
            const rules = config[field]
            // initialize fields
            this.fields[field] = {
                value: '',
                rules: [],
                message: '',
                error: false,
                touched: false,
            }
            if (rules.includes(',')) {
                const rulesArr = rules.split(',')
                for (const r in rulesArr) {
                    const rule = processRules(rulesArr[r])
                    this.fields[field].rules.push(rule)
                }
            } else {
                const rule = processRules(rules)
                this.fields[field].rules.push(rule)
            }
            this.$watch(`fields.${field}.value`, (value) => this.validate(field, value))
        }
    },
    // maybe put this somewhere else so we don't have to loop and validate all inputs
    // but for now it works
    paste: {
        ['@paste']() {
            for (const field in this.fields) {
                this.validate(field, this.fields[field].value)
            }
        },
    },
    must: {
        required: {
            validate(value) {
                if (!value || value === '') {
                    return 'required'
                }
            }
        },
        email: {
            validate(value) {
                if (!value.match(/^[\w-.]+@([\w-]+\.)+[\w-]{2,4}$/)) {
                    return 'invalid email'
                }
            }
        },
        minLength: {
            validate(value, length) {
                if (length && value.length < +length) {
                    return 'should be at least ' + length + ' characters'
                }
            }
        },
        maxLength: {
            validate(value, length) {
                if (length && value.length > +length) {
                    return 'should be at most ' + length + ' characters'
                }
            }
        },
        isNumber: {
            validate(value) {
                if (!isNaN(value)) {
                    return 'should be a number'
                }
            }
        },
        length: {
            validate(value, length) {
                if (length && value.length !== +length) {
                    return 'length should be ' + length
                }
            }
        }

    },
    isValid() {
        if (this.pristine || this.hasErrors) {
            return false
        }
        let allTouched = true
        for (const field in this.fields) {
            if (!this.fields[field].touched) {
                allTouched = false
                break
            }
        }
        return allTouched
    },
    validate(field, value) {
        const {rules} = this.fields[field]
        // set pristine to false if any field has a value
        this.pristine = false
        for (const r in rules) {
            const {rule, param} = rules[r]
            const validate = this.must[rule].validate
            if (validate) {
                let message = ''
                if (param === undefined) {
                    message = validate(value)
                } else {
                    message = validate(value, param)
                }
                if (message) {
                    this.fields[field].error = true
                    this.fields[field].message = message
                    this.fields[field].touched = true
                    break
                } else {
                    this.fields[field].error = false
                    this.fields[field].message = ''
                }
            }
        }
        // set global has errors at the form level
        let containsErrors = false
        for (const f in this.fields) {
            if (this.fields[f].error) {
                containsErrors = true
                break
            }
        }
        this.hasErrors = containsErrors
        this.valid = this.isValid()
    }
})

function processRules(rules) {
    if (rules.includes('=')) {
        const [rule, param] = rules.split('=')
        return {rule, param}
    } else {
        return {rule: rules, param: undefined}
    }
}