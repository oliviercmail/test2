<template>
  <b-card
    class="shadow-sm"
    header-bg-variant="white"
    footer-bg-variant="white"
  >
    <template #header>
      <h3 class="m-0">
        {{ $t('title') }}
      </h3>
    </template>

    <div>
      <div
        v-for="(colorInfo, key) in colorData"
        :key="key"
      >
        <label
          :for="key"
          class="text-primary"
        >{{ colorInfo.label }}</label>
        <div class="form-group row">
          <c-input-color-picker
            v-model="colorInfo.value"
            :data-test-id="`input-${key}-color`"
            width="36px"
            height="36px"
            show-color-code-text="true"
            class="px-3"
            :translations="colorTranslations"
          />
        </div>
      </div>
    </div>

    <template #footer>
      <c-submit-button
        :disabled="!canManage"
        :processing="processing"
        :success="success"
        class="float-right mt-2"
        @submit="onSubmit"
      />
    </template>
  </b-card>
</template>

<script>
import CSubmitButton from 'corteza-webapp-admin/src/components/CSubmitButton'
import { components } from '@cortezaproject/corteza-vue'
const { CInputColorPicker } = components

export default {
  name: 'CUIBrandingEditor',

  i18nOptions: {
    namespaces: 'ui.settings',
    keyPrefix: 'editor.corteza-studio',
  },

  components: {
    CInputColorPicker,
    CSubmitButton,
  },

  props: {
    settings: {
      type: Object,
      required: true,
    },

    processing: {
      type: Boolean,
      value: false,
    },

    success: {
      type: Boolean,
      value: false,
    },

    canManage: {
      type: Boolean,
      required: true,
    },
  },

  data () {
    return {
      colorData: {
        primary: {
          label: 'Primary Color',
          value: '',
        },
        secondary: {
          label: 'Secondary Color',
          value: '',
        },
      },
      colorTranslations: {
        modalTitle: this.$t('colorPicker'),
        saveBtnLabel: this.$t('general:label.saveAndClose'),
      },
    }
  },

  watch: {
    settings: {
      immediate: true,
      handler (settings) {
        const brandingVariables = settings['ui.branding.variables'] || []
        this.colorData.primary.value = this.brandColorValue(brandingVariables, 'primary')
        this.colorData.secondary.value = this.brandColorValue(brandingVariables, 'secondary')
      },
    },
  },

  methods: {
    onSubmit () {
      const brandVariables = Object.keys(this.colorData).map(key => {
        return `$${key}: ${this.colorData[key].value}`
      })

      this.$emit('submit', { 'ui.branding.variables': brandVariables })
    },

    brandColorValue (brandingVariables, variable) {
      const brandingVariable = brandingVariables.find(item => item.includes(`$${variable}:`))
      return brandingVariable.split(':')[1]
    },
  },
}
</script>

<style scoped lang="scss">
$border-color: #E4E9EF;
.input-group-text {
  padding: 0;
  background-color: var(--white);
}
.form-control:focus {
  border-color: $border-color;
}
</style>
