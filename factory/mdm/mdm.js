const products = require('./products.json').products;
const employees = require('./employees.json').employees;
const components = require('./components.json').components;

const db = { employees, components, products };

module.exports = () => ({
    employees: employees,
    components: components,
    products: products
});