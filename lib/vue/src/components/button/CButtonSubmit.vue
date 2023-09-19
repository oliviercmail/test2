<template>
  <b-button
    :data-test-id="cypressID"
    type="submit"
    :variant="variant"
    :disabled="disabled || processing || success"
    :size="buttonSize"
    :class="buttonClass"
    class="d-flex align-items-center justify-content-center text-center"
    @click.prevent="$emit('submit')"
  >
    <b-spinner
      v-if="processing"
      small
      :variant="iconVariant"
    />
    <font-awesome-icon
      v-else-if="success"
      :icon="['fas', 'check']"
      :class="iconVariant"
      class="text-white h3 mb-0"
    />
    <span
      v-else
    >
      {{ buttonText }}
    </span>
  </b-button>
</template>

<script>
export default {
  name: 'CSubmitButton',

  i18nOptions: {
    namespaces: 'admin',
  },

  props: {
    processing: {
      type: Boolean,
      value: false,
    },

    success: {
      type: Boolean,
      value: false,
    },

    disabled: {
      type: Boolean,
      value: false,
    },

    buttonClass: {
      type: String,
      default: '',
    },

    buttonText: {
      type: String,
      default: '',
    },

    buttonSize: {
      type: String,
      default: 'md',
    },

    variant: {
      type: String,
      default: 'primary',
    },

    iconVariant: {
      type: String,
      default: 'text-white',
    },

    cypressID: {
      type: String,
      default: 'button-submit',
    },
  },
}
</script>

<style>
.submit {
  min-width: 75px;
  min-height: 35px;
}
</style>
