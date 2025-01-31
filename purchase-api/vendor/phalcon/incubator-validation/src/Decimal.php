<?php

namespace Phalcon\Incubator\Validation;

use Phalcon\Filter\Validation;
use Phalcon\Messages\Message;
use Phalcon\Filter\Validation\AbstractValidator;
use Phalcon\Filter\Validation\Exception as ValidationException;
use Phalcon\Filter\Validation\ValidatorInterface;

/**
 * Phalcon\Validation\Validator\Decimal
 *
 * Validates that a value is a valid number in proper decimal format (negative and decimal numbers allowed).
 * Optionally, a specific number of digits can be checked too.
 *
 * Uses {@link http://www.php.net/manual/en/function.localeconv.php locale conversion} to allow decimal point to be
 * locale specific.
 *
 * <code>
 * use Phalcon\Validation\Validator\Decimal;
 *
 * $validator->add('number', new Decimal([
 *     'places'  => 2,
 *     'digit'   => 3,  // optional
 *     'point'   => ',' // optional. uses to override system decimal point
 *     'message' => 'Price must contain valid decimal value',
 * ]));
 * </code>
 */
class Decimal extends AbstractValidator implements ValidatorInterface
{
    /**
     * {@inheritdoc}
     *
     * @param Validation $validation
     * @param string $attribute
     *
     * @return bool
     * @throws ValidationException
     */
    public function validate(Validation $validation, $attribute): bool
    {
        $value = $validation->getValue($attribute);
        $field = $this->getOption('label');

        if (empty($field)) {
            $validation->getLabel($attribute);
        }

        if (false === $this->hasOption('places')) {
            throw new ValidationException(
                'A number of decimal places must be set'
            );
        }

        if ($this->hasOption('digits')) {
            // Specific number of digits
            $digits = '{' . ((int) $this->getOption('digits')) . '}';
        } else {
            // Any number of digits
            $digits = '+';
        }

        if ($this->hasOption('point')) {
            $decimal = $this->getOption('point');
        } else {
            // Get the decimal point for the current locale
            list($decimal) = array_values(
                localeconv()
            );
        }

        $result = (bool) preg_match(
            sprintf(
                '#^[+-]?[0-9]%s%s[0-9]{%d}$#',
                $digits,
                preg_quote($decimal),
                $this->getOption('places')
            ),
            $value
        );

        if (!$result) {
            $message = $this->getOption('message');

            $replacePairs = [
                ':field' => $field,
            ];

            if (empty($message)) {
                $message = ':field must contain valid decimal value';
            }

            $validation->appendMessage(
                new Message(
                    strtr($message, $replacePairs),
                    $attribute,
                    'Decimal'
                )
            );

            return false;
        }

        return true;
    }
}
