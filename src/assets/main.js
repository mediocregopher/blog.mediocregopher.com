// showModal will create the modal structure the first time it is called.
var modal, modalContent;
const showModal = function() {
  if (!modal) {
    // make the modal
    const modalClose = document.createElement('span');
    modalClose.id = 'modal-close';
    modalClose.innerHTML = '&times;';

    modalContent = document.createElement('div');
    modalContent.id = 'modal-content';

    const modalBody = document.createElement('div');
    modalBody.id = 'modal-body';
    modalBody.appendChild(modalContent);
    modalBody.appendChild(modalClose);

    modal = document.createElement('div');
    modal.id = 'modal';
    modal.appendChild(modalBody);

    // add the modal to the document
    document.getElementsByTagName('body')[0].appendChild(modal);

    // setup modal functionality
    modalClose.onclick = function() {
        modal.style.display = "none";
    }
  }

  modalContent.innerHTML = '';
  for (var i = 0; i < arguments.length; i++) {
    modalContent.appendChild(arguments[i]);
  }
  modal.style.display = "block";

  // When the user clicks anywhere outside of the modal, close it
  window.onclick = function(event) {
    if (event.target == modal) {
      modal.style.display = "none";
      window.onclick = undefined;
    }
  }
}

document.addEventListener("DOMContentLoaded", () => {
    console.log("DOM loaded");
})
