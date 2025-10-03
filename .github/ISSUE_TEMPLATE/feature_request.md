name: Feature Request
description: Suggest an idea for this project
title: "[Feature]: "
labels: [enhancement]
body:
  - type: textarea
    id: problem
    attributes: 
      label: Is your feature request related to a problem?
      description: Describe the frustration or gap this feature solves.
      placeholder: "I'm frustrated when..."
    validations:
      required: false

  - type: textarea
    id: solution
    attributes:
      label: Describe the solution you'd like
      description: Be as clear as possible.
      placeholder: "I’d like the system to..."
    validations:
      required: true