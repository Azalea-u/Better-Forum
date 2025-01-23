import Navbar from './components/Navbar.js';
import Footer from './components/Footer.js';
import Home from './pages/Home.js';
import Login from './pages/Login.js';
import Register from './pages/Register.js';

const App = {
  init() {
    document.getElementById('app').innerHTML = `
      <header id="navbar"></header>
      <main id="content"></main>
      <footer id="footer"></footer>
    `;
    this.checkAuth();
    Navbar.render();
    Footer.render();
    this.setupListeners();
  },

  checkAuth() {
    const isLoggedIn = sessionStorage.getItem('isLoggedIn');
    if (!isLoggedIn) {
      this.route('login'); // Redirect to login if not authenticated
    } else {
      this.route('home'); // Redirect to home if logged in
    }
  },

  route(page) {
    const content = document.getElementById('content');
    switch(page) {
      case 'login':
        Login.render(content);
        break;
      case 'register':
        Register.render(content);
        break;
      default:
        Home.render(content); // Default: Home
        break;
    }
  },

  setupListeners() {
    document.addEventListener('navigate', (e) => {
      this.route(e.detail.page);
    });
  },
};

export default App;
