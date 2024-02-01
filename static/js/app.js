// app.js
//
(function () {
  "use strict";
  window.addEventListener(
    "load",
    function () {
      // Fetch all the forms we want to apply custom Bootstrap validation styles to
      let forms = document.getElementsByClassName("needs-validation");
      // Loop over them and prevent submission
      Array.prototype.filter.call(forms, function (form) {
        form.addEventListener(
          "submit",
          function (event) {
            if (form.checkValidity() === false) {
              event.preventDefault();
              event.stopPropagation();
            }
            form.classList.add("was-validated");
          },
          false,
        );
      });
    },
    false,
  );
})();

function createAutoDismissNotification(msg, msgType) {
  // Set the modal title
  var modalTitle = document.getElementById("notificationModalLabel");
  modalTitle.innerHTML = msgType;

  // Set the modal body content
  var modalBody = document.getElementById("notificationModalBody");
  modalBody.innerHTML = msg;

  // Add the 'show' class to the modal to make it visible
  var notificationModal = new bootstrap.Modal(
    document.getElementById("notificationModal"),
  );
  notificationModal.show();

  // Close the modal after a specified duration (e.g., 10 seconds)
  setTimeout(function () {
    // Hide the modal using Bootstrap's hide method
    notificationModal.hide();
  }, 10000);
}
