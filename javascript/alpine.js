import Alpine from 'alpinejs'
import validation from "./validation.js";

window.Alpine = Alpine
Alpine.data('validation', validation)
Alpine.start()