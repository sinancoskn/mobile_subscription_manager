<?php

namespace Phalcon\Incubator\Validation;

use Phalcon\Filter\Validation;
use Phalcon\Messages\Message;
use Phalcon\Filter\Validation\AbstractValidator;
use Phalcon\Filter\Validation\ValidatorInterface;

class AlphaNamesValidator extends AbstractValidator implements ValidatorInterface
{
    /**
     * Executes the validation. Allowed options:
     * 'numbers' : allow numbers;
     * 'min' : input value must not be shorter than it;
     * 'max' : input value must not be longer than it.
     *
     * @param  Validation $validator
     * @param  string $attribute
     *
     * @return boolean
     */
    public function validate(Validation $validator, $attribute): bool
    {
        $value = $validator->getValue($attribute);

        $numbers = (bool) $this->getOption('numbers');
        $numbers = $numbers ? '0-9' : '';

        if (!preg_match('/^([-\p{L}' . $numbers . '\'_\s])+$/u', $value)) {
            $message = $this->getOption('message');

            if (!$message) {
                if ($numbers) {
                    $message = 'The value can contain only alphanumeric, menus, apostrophe, underscore and '
                        . 'white space characters';
                } else {
                    $message = 'The value can contain only alphabetic, menus, apostrophe, underscore and '
                        . 'white space characters';
                }
            }

            $validator->appendMessage(
                new Message(
                    $message,
                    $attribute,
                    'AlphaNames'
                )
            );
        }

        if ($min = (int)$this->getOption('min')) {
            if (strlen($value) < $min) {
                $messageMin = $this->getOption(
                    'messageMinimum',
                    'The value must contain at least ' . $min . ' characters.'
                );

                $validator->appendMessage(
                    new Message(
                        $messageMin,
                        $attribute,
                        'AlphaNames'
                    )
                );
            }
        }

        if ($max = (int) $this->getOption('max')) {
            if (strlen($value) > $max) {
                $messageMax = $this->getOption(
                    'messageMaximum',
                    'The value can contain maximum ' . $max . ' characters.'
                );

                $validator->appendMessage(
                    new Message(
                        $messageMax,
                        $attribute,
                        'AlphaNames'
                    )
                );

                return false;
            }
        }

        if (count($validator->getMessages())) {
            return false;
        }

        return true;
    }
}
