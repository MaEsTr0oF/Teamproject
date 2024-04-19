// Подключение функционала "Чертогов Фрилансера"
//import { isMobile } from "./functions.js";
// Подключение списка активных модулей
//import { flsModules } from "./modules.js";

document.addEventListener("click", documentActions);

function documentActions(e) {
	const targetElement = e.target;
	if (targetElement.closest('[data-parent]')) {
		const subMenuId = targetElement.dataset.parent ? targetElement.dataset.parent : null;
		const subMenu = document.querySelector(`[data-submenu="${subMenuId}"]`);
		const catalogMenu = document.querySelector(`.catalog-header`);
		if (subMenu) {
			catalogMenu.classList.toggle('_sub-menu-show');
			targetElement.classList.toggle('_sub-menu-active');
			subMenu.classList.toggle('_sub-menu-open');
		}
		else {
			console.log(`Ой ой ,тут нет ничего`)
		}
		e.preventDefault()
	}
}